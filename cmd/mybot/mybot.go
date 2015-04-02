package main

import (
	"flag"
	"fmt"
	"github.com/cenkalti/rpc2"
	irc "github.com/fluffle/goirc/client"
	bot "github.com/jebjerg/go-bot/bot"
	cfg "github.com/jebjerg/go-bot/bot/config"
	"net"
)

type bot_conf struct {
	IRCHost       string `json:"irc_addr"`
	ListenAddr    string `json:"listen_addr"`
	LuaScriptPath string `json:"lua_scripts"`
}

func main() {
	var rpc_support, lua_support bool
	flag.BoolVar(&rpc_support, "rpc", true, "enable RPC support")
	flag.BoolVar(&lua_support, "lua", true, "enable Lua support")
	flag.Parse()
	conf := &bot_conf{}
	cfg.NewConfig(conf, "mybot.json")

	ir := bot.NewClient("bot")
	ir.Client.Config().SSL = false
	ir.Client.Config().Server = conf.IRCHost
	ir.Client.Config().NewNick = func(nick string) string { return nick + "_" }

	ir.Client.HandleFunc("connected", func(conn *irc.Conn, line *irc.Line) {
		fmt.Println("online")
		conn.Join("#channel")
	})

	quit := make(chan bool)
	ir.Client.HandleFunc("disconnected", func(conn *irc.Conn, line *irc.Line) {
		quit <- true
	})
	var connected bool
	ir.Connect(nil, &connected)

	ir.Client.HandleFunc("privmsg", func(conn *irc.Conn, line *irc.Line) {
		if lua_support {
			ir.PrivMsgLua(line.Args[0], line.Args[1])
		}
		if rpc_support {
			for c, _ := range ir.Listeners {
				if err := c.Call("privmsg", line, nil); err != nil {
					delete(ir.Listeners, c)
				}
			}
		}
	})

	if lua_support {
		if conf.LuaScriptPath != "" {
			ir.LuaScriptPath = conf.LuaScriptPath
		}
		ir.InitLua()
	}
	if rpc_support {
		l, e := net.Listen("tcp", conf.ListenAddr)
		defer l.Close()
		if e != nil {
			panic(e)
		}
		srv := rpc2.NewServer()
		srv.Handle("register", ir.Register)
		srv.Handle("privmsg", ir.Announce)
		srv.Handle("join", ir.Join)
		go srv.Accept(l)
		fmt.Println("RPC ready")
	}

	fmt.Println("IRC process running")
	<-quit
}
