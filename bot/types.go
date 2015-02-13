package bot

import (
	"github.com/cenkalti/rpc2"
	irc "github.com/fluffle/goirc/client"
)

type VoidArgs struct{}

type PrivMsgArgs struct {
	Target, Text string
}

type IRCRPC struct {
	Client    *irc.Conn
	Listeners map[*rpc2.Client]bool
}
