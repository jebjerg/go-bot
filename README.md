# go-bot

## Go-what?
It's a small irc bot written in Go, as an exercise, with bi-directional RPC and Lua interfaces

## Features
* Written in [Go](http://golang.org)
* Minimalistic and easy to use
* Easy to extend with any language (via RPC) or [Lua 5.1](http://www.lua.org/manual/5.1/)

See [`cmd/mybot`](cmd/mybot) for an example bot.
See [`cmd/boilerplate_rpc`](cmd/boilerplate_rpc) for a simple boilerplate RPC plugin.
See [`cmd/mybot/boilerplate.lua`](cmd/mybot/boilerplate.lua) for a simple Lua plugin.

See [twitter](https://github.com/jebjerg/bot-twitter) for an example of a RPC plugin.
See [uptime](https://github.com/jebjerg/bot-uptime) for an example of a Lua plugin.

More plugins in [my repos](https://github.com/jebjerg?tab=repositories).

## Getting started
### mybot
Assuming you have `go` installed:
```
# go get github.com/jebjerg/go-bot/cmd/mybot
# mkdir -p mybot/scripts
# curl -LO https://github.com/jebjerg/go-bot/raw/master/cmd/mybot/mybot.json -o mybot/mybot.json
# curl -LO https://github.com/jebjerg/bot-uptime/raw/master/uptime.lua -o mybot/uptime.lua
```

```
# mybot -h
Usage of mybot:
  -lua=true: enable Lua support
  -rpc=true: enable RPC support
```
Put scripts in `./` or specify path in `mybot.json` (e.g. `"lua_scripts": "scripts"` to use `mybot/scripts/uptime.lua`).
Reload scripts by issuing a `SIGKILL`, `SIGUSR1` or `SIGUSR2`, e.g. `pkill -HUP mybot`.

```

```

If you need help getting started with `go`, follow these [two steps to install](https://golang.org/doc/install#install)

### dev
For Go dev'ing, the regular steps are as expected: `go get github.com/jebjerg/go-bot`
Use [`cmd/mybot`](cmd/mybot) as a starting point.

## Todo:
* Complete API (for RPC and Lua)
* Documentation
* RPC connectivity bug luring
