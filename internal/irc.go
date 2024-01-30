package bridge

import (
	"fmt"
	"os"

	"github.com/whyrusleeping/hellabot"
	"gopkg.in/sorcix/irc.v2"
)

type Irc struct {
	bot     *hbot.Bot
	discord *Discord
}

func NewIrcBridge() Irc {
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

	return Irc{
		bot,
		nil,
	}
}

func (i *Irc) setDiscord(discord *Discord) {
	i.discord = discord
}

func (i *Irc) Setup() {
	i.bot.AddTrigger(hbot.Trigger{
		Condition: func(bot *hbot.Bot, msg *hbot.Message) bool {
			return msg.Command == irc.PRIVMSG
		},
		Action: func(bot *hbot.Bot, msg *hbot.Message) bool {
			fmt.Println(msg)
			i.discord.sendMessage(msg.Content)

			return true
		},
	})
}

func (i *Irc) Run() {
	fmt.Println("irc bot is up")
	i.bot.Run()
}

func (i *Irc) sendMessage(msg string) {
	i.bot.Msg("#dev", msg)
}
