package main

import (
	bot "./bot"
	"fmt"
	"github.com/cenkalti/rpc2"
	irc "github.com/fluffle/goirc/client"
	"net"
)

func main() {
	ir := bot.NewClient("bot")
	ir.Client.Config().SSL = false
	ir.Client.Config().Server = "127.0.0.1:6667" // "irc.freenode.net:7000"
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
		var reply bool
		for c, _ := range ir.Listeners {
			if err := c.Call("privmsg", line, &reply); err != nil {
				delete(ir.Listeners, c)
			}
		}
	})

	// RPC
	l, e := net.Listen("tcp", ":1234")
	defer l.Close()
	if e != nil {
		panic(e)
	}

	srv := rpc2.NewServer()
	srv.Handle("register", ir.Register)
	srv.Handle("privmsg", ir.Announce)
	srv.Handle("join", ir.Join)
	go srv.Accept(l)

	fmt.Println("IRC process running, ready for RPC")
	<-quit
}
