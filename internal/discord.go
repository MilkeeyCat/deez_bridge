package bridge

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var channelId string
var guildId string

type Discord struct {
	bot *discordgo.Session
	irc *Irc
}

func NewDiscordBridge() Discord {
	guildId = os.Getenv("DISCORD_BOT_GUILD_ID")
	channelId = os.Getenv("DISCORD_BOT_CHANNEL_ID")
	token := os.Getenv("DISCORD_BOT_TOKEN")

	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	return Discord{
		bot: bot,
		irc: nil,
	}
}

func (d *Discord) setIrc(irc *Irc) {
	d.irc = irc
}

func (d *Discord) Setup() {
	d.bot.AddHandler(func(discord *discordgo.Session, message *discordgo.MessageCreate) {
		if discord.State.User.ID != message.Author.ID && message.ChannelID == channelId {
			d.irc.sendMessage(message.Author.Username, message.Content)
		}
	})
}

func (b *Discord) Run() {
	err := b.bot.Open()
	if err != nil {
		panic(err)
	}

	defer b.bot.Close()

	fmt.Println("deezcord bot is up")

	// is there a way to make it prettier?
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func (d *Discord) sendMessage(author, msg string) {
	_, err := d.bot.ChannelMessageSend(channelId, fmt.Sprintf("<%s> %s", author, msg))
	if err != nil {
		fmt.Println(err)
		return
	}
}
