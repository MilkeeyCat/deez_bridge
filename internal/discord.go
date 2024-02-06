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
	bot.AddHandler(discord.onEdit)
	bot.AddHandler(discord.onReactionAdd)
	bot.AddHandler(discord.onReactionRemove)

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

func (d *Discord) replyToMessage(username string, author string, content string, offset int32) {
	message := d.bridge.messages.findByOffset(username, uint32(offset))

	if message != nil {
		message, err := d.bot.ChannelMessageSendReply(channelId, fmt.Sprintf("<%s> %s", author, content), &discordgo.MessageReference{
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

func (d *Discord) onEdit(session *discordgo.Session, message *discordgo.MessageUpdate) {
	author := message.Author.Username
	content := message.Content
	offset := d.bridge.messages.find(author, message.ID)

	d.bridge.messages.update(message.ID, author, content)
	d.bridge.irc.sendMessage(fmt.Sprintf("<%s~%d> %s", author, offset, content))
}

func (d *Discord) onReactionAdd(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	emojiName := reaction.Emoji.Name
	from := reaction.Member.User.Username
	to := ""
	msg, err := session.ChannelMessage(channelId, reaction.MessageID)
	if err != nil {
		fmt.Printf("failed to get message: %s", err)
	}

	if session.State.User.ID == msg.Author.ID {
		to = MessageAuthor(msg.Content)
	} else {
		to = msg.Author.Username
	}

	messageId := d.bridge.messages.find(to, reaction.MessageID)
	content := fmt.Sprintf("%s reacted with %s to %s~%d", from, emojiName, to, messageId)

	d.bridge.irc.sendMessage(content)
}

func (d *Discord) onReactionRemove(session *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	emojiName := reaction.Emoji.Name
	user, err := session.User(reaction.UserID)
	if err != nil {
		fmt.Println(err)
	}
	from := user.Username
	to := ""
	msg, err := session.ChannelMessage(channelId, reaction.MessageID)
	if err != nil {
		fmt.Printf("failed to get message: %s", err)
	}

	if session.State.User.ID == msg.Author.ID {
		to = MessageAuthor(msg.Content)
	} else {
		to = msg.Author.Username
	}

	messageId := d.bridge.messages.find(to, reaction.MessageID)
	content := fmt.Sprintf("%s removed reaction %s from %s~%d", from, emojiName, to, messageId)

	d.bridge.irc.sendMessage(content)
}

func (d *Discord) deleteMessage(name string, offset int) {
	message := d.bridge.messages.findByOffset(name, uint32(offset))
	if message == nil {
		fmt.Println("failed to find message")
		return
	}

	d.bot.ChannelMessageDelete(channelId, message.messageId)
	d.bridge.messages.delete(name, message.messageId)
}
