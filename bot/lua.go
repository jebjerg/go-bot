package bot

import (
	"crypto/tls"
	"fmt"
	lua "github.com/yuin/gopher-lua"
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
		"hook":           c.HookLua,
		"connect":        c.ConnectLua,
		"connected":      c.ConnectedLua,
		"set_ssl":        c.SetSSLLua,
		"set_ssl_verify": c.SetSSLVerifyLua,
		"quit":           c.QuitLua,
		"raw":            c.RawIRCLua,
		"privmsg":        c.SendPrivMsgLua,
		"nick":           c.NickLua,
		"join":           c.JoinLua,
		"part":           c.PartLua,
		"topic":          c.TopicLua,
	}
}

func (c *IRCRPC) LoadScripts() {
	c.LuaStates = make(map[*lua.LState]map[string][]*lua.LFunction)
	files, err := filepath.Glob(path.Join(c.LuaScriptPath, "*.lua"))
	if err != nil {
		return
	}
	for _, file := range files {
		go func(file string) {
			state := lua.NewState()
			c.BootstrapState(state)
			c.LuaStates[state] = make(map[string][]*lua.LFunction)
			if err := state.DoFile(file); err != nil {
				c.LuaStates[state] = nil
				state = nil
				fmt.Println("WARNING:", file, "failed:", err)
			}
			fmt.Println(file, "loaded")
		}(file)
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

func (c *IRCRPC) ConnectLua(L *lua.LState) int {
	host := L.ToString(1)
	pass := L.ToString(2)
	c.Client.ConnectTo(host, pass)
	return 0
}

func (c *IRCRPC) ConnectedLua(L *lua.LState) int {
	L.Push(lua.LBool(c.Client.Connected()))
	return 1
}

func (c *IRCRPC) SetSSLLua(L *lua.LState) int {
	c.Client.Config().SSL = L.ToBool(1)
	if c.Client.Config().SSL == true && c.Client.Config().SSLConfig == nil {
		c.Client.Config().SSLConfig = &tls.Config{}
	}
	return 0
}

func (c *IRCRPC) SetSSLVerifyLua(L *lua.LState) int {
	if c.Client.Config().SSLConfig == nil {
		c.Client.Config().SSLConfig = &tls.Config{}
	}
	c.Client.Config().SSLConfig.InsecureSkipVerify = L.ToBool(1)
	return 0
}

func (c *IRCRPC) QuitLua(L *lua.LState) int {
	msg := L.ToString(1)
	c.Client.Quit(msg)
	return 0
}

func (c *IRCRPC) RawIRCLua(L *lua.LState) int {
	raw := L.ToString(1)
	c.Client.Raw(raw)
	return 0
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

func (c *IRCRPC) NickLua(L *lua.LState) int {
	new_nick := L.ToString(1)
	c.Client.Nick(new_nick)
	return 0
}

func (c *IRCRPC) JoinLua(L *lua.LState) int {
	target := L.ToString(1)
	c.Client.Join(target)
	return 0
}

func (c *IRCRPC) PartLua(L *lua.LState) int {
	target := L.ToString(1)
	msg := L.ToString(2)
	c.Client.Part(target, msg)
	return 0
}

func (c *IRCRPC) TopicLua(L *lua.LState) int {
	target := L.ToString(1)
	topic := L.ToString(2)
	c.Client.Topic(target, topic)
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
