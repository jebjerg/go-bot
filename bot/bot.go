package bot

import (
	"github.com/cenkalti/rpc2"
	irc "github.com/fluffle/goirc/client"
)

func NewClient(nick string) *IRCRPC {
	lm := make(map[*rpc2.Client]bool)
	return &IRCRPC{Client: irc.SimpleClient(nick, nick), Listeners: lm, LuaScriptPath: "."}
}
