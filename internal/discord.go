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
	bot      *discordgo.Session
	irc      *Irc
	messages *MessagesMap
}

func NewDiscordBridge(messages *MessagesMap) Discord {
	guildId = os.Getenv("DISCORD_BOT_GUILD_ID")
	channelId = os.Getenv("DISCORD_BOT_CHANNEL_ID")
	token := os.Getenv("DISCORD_BOT_TOKEN")

	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	return Discord{
		bot,
		nil,
		messages,
	}
}

func (d *Discord) setIrc(irc *Irc) {
	d.irc = irc
}

func (d *Discord) Setup() {
	d.bot.AddHandler(func(discord *discordgo.Session, message *discordgo.MessageCreate) {
		author := message.Author.Username
		content := message.Content

		if discord.State.User.ID != message.Author.ID && message.ChannelID == channelId && message.Type == discordgo.MessageTypeDefault {
			formatedContent := fmt.Sprintf("<%s> %s", author, content)

			d.messages.push(author, content)
			d.irc.sendMessage(author, formatedContent)
		} else if message.Type == discordgo.MessageTypeReply {
			reply := message.ReferencedMessage
			to := ""
			replyContent := reply.Content

			if discord.State.User.ID == reply.Author.ID {
				to = MessageAuthor(replyContent)
			} else {
				to = reply.Author.Username
			}

			d.messages.push(author, content)
			i := d.messages.find(to, MessageContent(replyContent))
			d.irc.sendMessage(to, fmt.Sprintf("<%s ^%d %s> %s", author, i, to, content))

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
	_, err := d.bot.ChannelMessageSend(channelId, msg)
	if err != nil {
		fmt.Println(err)
		return
	}
}
