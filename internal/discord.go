package bridge

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

var channelId string
var guildId string

type Discord struct {
	bot    *discordgo.Session
	bridge *Bridge
}

func NewDiscordBridge(bridge *Bridge) *Discord {
	guildId = os.Getenv("DISCORD_BOT_GUILD_ID")
	channelId = os.Getenv("DISCORD_BOT_CHANNEL_ID")
	token := os.Getenv("DISCORD_BOT_TOKEN")

	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	discord := &Discord{
		bot,
		bridge,
	}

	bot.AddHandler(discord.onMessage)

	return discord
}

func (d *Discord) Open() {
	err := d.bot.Open()
	if err != nil {
		panic(err)
	}
	fmt.Println("discord connection established")
}

func (d *Discord) Close() {
	err := d.bot.Close()
	if err != nil {
		panic(err)
	}
}

func (d *Discord) sendMessage(message string) {
	msg, err := d.bot.ChannelMessageSend(channelId, message)
	if err != nil {
		fmt.Println("error occurred during sending message: %w", err)
	}

	d.bridge.messages.push(MessageAuthor(message), Message{
		content:   MessageContent(message),
		messageId: msg.ID,
	})
}

func (d *Discord) onMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if session.State.User.ID != message.Author.ID && message.ChannelID == channelId && message.Type == discordgo.MessageTypeDefault {
		author := message.Author.Username
		content := message.Content

		msg := fmt.Sprintf("<%s> %s", author, content)

		d.bridge.irc.sendMessage(msg)
		d.bridge.messages.push(author, Message{
			content:   content,
			messageId: message.Message.ID,
		})

	} else if session.State.User.ID != message.Author.ID && message.ChannelID == channelId && message.Type == discordgo.MessageTypeReply {
		d.onReply(session, message)
	}
}

func (d *Discord) onReply(session *discordgo.Session, message *discordgo.MessageCreate) {
	repliedMessage := message.ReferencedMessage
	from := message.Author.Username
	content := message.Content
	to := ""

	if session.State.User.ID == repliedMessage.Author.ID {
		to = MessageAuthor(repliedMessage.Content)
	} else {
		to = repliedMessage.Author.Username
	}

	messageId := d.bridge.messages.find(to, repliedMessage.ID)

	msg := fmt.Sprintf("<%s %s~%d> %s", from, to, messageId, content)

	d.bridge.irc.sendMessage(msg)
	d.bridge.messages.push(from, Message{
		content:   content,
		messageId: message.Message.ID,
	})
}

func (d *Discord) replyToMessage(username string, content string, offset int32) {
	message := d.bridge.messages.findByOffset(username, uint32(offset))

	if message != nil {
		message, err := d.bot.ChannelMessageSendReply(channelId, content, &discordgo.MessageReference{
			MessageID: message.messageId,
			ChannelID: channelId,
			GuildID:   guildId,
		})

		if err != nil {
			fmt.Println("failed to reply to a message: %w", err)
		}

		d.bridge.messages.push(username, Message{
			content:   content,
			messageId: message.ID,
		})
	}
}
