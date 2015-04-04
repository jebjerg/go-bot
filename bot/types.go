package bot

import (
	"github.com/cenkalti/rpc2"
	irc "github.com/fluffle/goirc/client"
	lua "github.com/yuin/gopher-lua"
)

type VoidArgs struct{}

type PrivMsgArgs struct {
	Target, Text string
}

type IRCRPC struct {
	Client        *irc.Conn
	LuaStates     map[*lua.LState]map[string][]*lua.LFunction
	LuaScriptPath string
	Listeners     map[*rpc2.Client]bool
}
