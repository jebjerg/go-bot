package main

import (
	"fmt"
	"github.com/cenkalti/rpc2"
	irc "github.com/fluffle/goirc/client"
	bot "github.com/jebjerg/go-bot/bot"
	cfg "github.com/jebjerg/go-bot/bot/config"
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
	BotHost  string   `json:"bot_host"`
}

func main() {
	conf := &boilerplate_conf{}
	cfg.NewConfig(conf, "boilerplate.json")

	// RPC
	conn, _ := net.Dial("tcp", conf.BotHost)
	c := rpc2.NewClient(conn)
	go c.Run()
	// register privmsg
	c.Handle("privmsg", OnPrivMsg)
	var reply bool
	c.Call("register", struct{}{}, &reply)

	for _, channel := range conf.Channels {
		fmt.Println("joining", channel)
		c.Call("join", channel, nil)
	}

	// daemon
	forever := make(chan bool)
	<-forever
}
