package discord

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

type BridgeDiscord struct {
	bot *discordgo.Session
}

func NewDiscordBridge() BridgeDiscord {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	channelId := os.Getenv("DISCORD_BOT_CHANNEL_ID")

	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	bot.AddHandler(func(discord *discordgo.Session, message *discordgo.MessageCreate) {
		if discord.State.User.ID != message.Author.ID {
			_, err := discord.ChannelMessageSend(channelId, "hi :-)")
			if err != nil {
				panic(err)
			}
		}
	})

	return BridgeDiscord{
		bot,
	}
}

func (bridge BridgeDiscord) Run() {
	err := bridge.bot.Open()
	if err != nil {
		panic(err)
	}
	defer bridge.bot.Close()

	fmt.Println("deezcord bot is up")

	// is there a way to make it prettier?
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
