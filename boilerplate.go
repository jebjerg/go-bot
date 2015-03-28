package main

import (
	bot "./bot"
	cfg "./bot/config"
	"fmt"
	"github.com/cenkalti/rpc2"
	irc "github.com/fluffle/goirc/client"
	"net"
)

func OnPrivMsg(client *rpc2.Client, line *irc.Line, reply *bool) error {
	channel, text := line.Args[0], line.Args[1]
	if text == "!help" {
		msg := "I'm not sure I can"
		client.Call("privmsg", &bot.PrivMsgArgs{channel, msg}, nil)
	}
	return nil
}

type boilerplate_conf struct {
	Channels []string `json:"channels"`
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

	conf := &boilerplate_conf{}
	cfg.NewConfig(conf, "boilerplate.json")
	for _, channel := range conf.Channels {
		fmt.Println("joining", channel)
		c.Call("join", channel, nil)
	}

	// daemon
	forever := make(chan bool)
	<-forever
}
