package main

import (
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"net"
	"net/http"
	"net/rpc"
)

type VoidArgs struct{}

type PrivMsgArgs struct {
	Target, Text string
}

type IRCRPC struct {
	Client *irc.Conn
}

func NewClient(nick string) *IRCRPC {
	return &IRCRPC{irc.SimpleClient(nick, nick)}
}

func (c *IRCRPC) Connect(args *VoidArgs, reply *bool) error {
	if err := c.Client.Connect(); err != nil {
		*reply = false
		return err
	}
	*reply = true
	return nil
}

func (c *IRCRPC) Disconnect(args *VoidArgs, reply *string) error {
	c.Client.Quit("bye bye")
	return nil
}

func (c *IRCRPC) Announce(msg *PrivMsgArgs, reply *bool) error {
	c.Client.Privmsg(msg.Target, msg.Text)
	*reply = true
	return nil
}

func main() {
	ir := NewClient("bot")
	ir.Client.Config().SSL = false
	ir.Client.Config().Server = "127.0.0.1:6667" // "irc.freenode.net:7000"
	ir.Client.Config().NewNick = func(nick string) string { return nick + "_" }

	ir.Client.HandleFunc("connected", func(conn *irc.Conn, line *irc.Line) {
		fmt.Println("c0nnected")
		conn.Join("#channel")
	})

	quit := make(chan bool)
	ir.Client.HandleFunc("disconnected", func(conn *irc.Conn, line *irc.Line) {
		quit <- true
	})
	var connected bool
	ir.Connect(nil, &connected)

	rpc.Register(ir)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	defer l.Close()
	if e != nil {
		panic(e)
	}
	go http.Serve(l, nil)
	fmt.Println("IRC process running, ready for RPC")
	<-quit
}
