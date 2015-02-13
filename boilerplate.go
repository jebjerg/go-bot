package main

import (
	"github.com/cenkalti/rpc2"
	irc "github.com/fluffle/goirc/client"
	"net"
	"strings"
)

// will be moved to main bot's package, along with irc.Line
type PrivMsg struct {
	Target, Text string
}

func OnPrivMsg(client *rpc2.Client, line *irc.Line, reply *bool) error {
	channel, text := line.Args[0], line.Args[1]
	if text == "!help" {
		msg := "I'm not sure I can"
		client.Call("privmsg", &PrivMsg{channel, msg}, &reply)
	}
	return nil
}

func main() {
	// RPC
	conn, _ := net.Dial("tcp", "localhost:1234")
	c := rpc2.NewClient(conn)
	go c.Run()
	// register privmsg
	c.Handle("privmsg", OnPrivMsg)
	var reply bool
	c.Call("register", struct{}{}, &reply)

	// daemon
	forever := make(chan bool)
	<-forever
}
