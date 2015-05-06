-- define event handlers, callbacks
function echo(channel, text)
    bot.privmsg(channel, text)
end

function log_privmsg(channel, text)
    print("logging: "..channel..": "..text)
end

-- some lua code, do your stuff
local name = "boilerplate"
print(name.." loaded, whoo!")

-- bot preloads and sets 'bot' table of functions representing the lua api
bot.hook("PRIVMSG", "echo")
bot.hook("PRIVMSG", "log_privmsg")
