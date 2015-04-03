package bot

import (
	"fmt"
	lua "github.com/jebjerg/gopher-lua"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"
)

func (c *IRCRPC) BootstrapState(L *lua.LState) int {
	// Alternatively, use PreloadModule as before, and do local bot = require("bot") in all scripts.
	bot := L.SetFuncs(L.NewTable(), c.LuaAPI())
	// L.SetField(bot, "version", lua.LString("v0.1"))
	L.Push(bot)
	L.SetGlobal("bot", bot)
	return 1
}

func (c *IRCRPC) LuaAPI() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"hook":    c.HookLua,
		"privmsg": c.SendPrivMsgLua,
	}
}

func (c *IRCRPC) LoadScripts() {
	c.LuaStates = make(map[*lua.LState]map[string][]*lua.LFunction)
	files, err := filepath.Glob(path.Join(c.LuaScriptPath, "*.lua"))
	if err != nil {
		return
	}
	for _, file := range files {
		state := lua.NewState()
		c.BootstrapState(state)
		c.LuaStates[state] = make(map[string][]*lua.LFunction)
		if err := state.DoFile(file); err != nil {
			c.LuaStates[state] = nil
			state = nil
			fmt.Println("WARNING:", file, "failed:", err)
			continue
		}
		fmt.Println(file, "loaded")
	}
}

func (c *IRCRPC) Cleanup() {
	c.LuaStates = nil
}

func (c *IRCRPC) HookLua(L *lua.LState) int {
	event := L.ToString(1)
	fn := L.ToFunction(2)
	if event == "" || fn == nil {
		return 1
	}
	if c.LuaStates[L][event] != nil {
		c.LuaStates[L][event] = append(c.LuaStates[L][event], fn)
	} else {
		funcs := []*lua.LFunction{fn}
		c.LuaStates[L][event] = funcs
	}
	return 0
}

func (c *IRCRPC) Trigger(event string, args ...lua.LValue) error {
	for state, fnmap := range c.LuaStates {
		if fns := fnmap[event]; fns != nil {
			for _, fn := range fns {
				state.CallByParam(lua.P{
					Fn:      fn,
					NRet:    0,
					Protect: true,
				}, args...)
			}
		}
	}
	return nil
}

func (c *IRCRPC) PrivMsgLua(channel, text string) error {
	return c.Trigger("PRIVMSG", lua.LString(channel), lua.LString(text))
}

func (c *IRCRPC) SendPrivMsg(target, text string) {
	c.Client.Privmsg(target, text)
}

func (c *IRCRPC) SendPrivMsgLua(L *lua.LState) int {
	target := L.ToString(1)
	text := L.ToString(2)
	c.SendPrivMsg(target, text)
	return 0
}

func (c *IRCRPC) InitLua() {
	fmt.Println("Lua engine ready")
	c.LoadScripts()
	go func() {
		reload_sig := make(chan os.Signal, 1)
		signal.Notify(reload_sig, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP)
		for {
			<-reload_sig
			fmt.Println("reloading scripts")
			c.Cleanup()
			c.LoadScripts()
		}
	}()
}
