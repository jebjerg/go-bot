package bot

import (
	lua "github.com/aarzilli/golua/lua"
	"github.com/cenkalti/rpc2"
	irc "github.com/fluffle/goirc/client"
)

type VoidArgs struct{}

type PrivMsgArgs struct {
	Target, Text string
}

type IRCRPC struct {
	Client        *irc.Conn
	LuaStates     map[*lua.State]map[string][]string
	LuaScriptPath string
	Listeners     map[*rpc2.Client]bool
}
