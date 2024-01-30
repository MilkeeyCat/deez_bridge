package irc

import (
	"fmt"
	"os"

	"github.com/whyrusleeping/hellabot"
)

type BridgeIrc struct {
	bot *hbot.Bot
}

func NewIrcBridge() BridgeIrc {
	host := os.Getenv("IRC_SERVER_HOST")
	port := os.Getenv("IRC_SERVER_PORT")
	nickname := os.Getenv("IRC_NICKNAME")
	channel := os.Getenv("IRC_CHANNEL")

	bot, err := hbot.NewBot(fmt.Sprintf("%s:%s", host, port), nickname, func(bot *hbot.Bot) {
		bot.Channels = []string{channel}
	})
	if err != nil {
		panic(err)
	}

	bot.AddTrigger(hbot.Trigger{
		Condition: func(bot *hbot.Bot, msg *hbot.Message) bool {
			return true
		},
		Action: func(bot *hbot.Bot, msg *hbot.Message) bool {
            bot.Reply(msg, "reply :)")

			return true
		},
	})

	return BridgeIrc{
		bot: bot,
	}
}

func (bridge BridgeIrc) Run() {
    fmt.Println("irc bot is up")
	bridge.bot.Run()
}
