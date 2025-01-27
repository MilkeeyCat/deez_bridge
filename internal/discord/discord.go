package discord

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/MilkeeyCat/deez_bridge/internal/logger"
	"github.com/MilkeeyCat/deez_bridge/internal/message"
	"github.com/bwmarrin/discordgo"
)

var (
	channelId string
	guildId   string
)

type Discord struct {
	bot           *discordgo.Session
	nickMemberMap map[string]*discordgo.Member
	messages      Messages
	message       chan message.Message
}

func NewDiscord(messageChan chan message.Message) Discord {
	guildId = os.Getenv("DISCORD_BOT_GUILD_ID")
	channelId = os.Getenv("DISCORD_BOT_CHANNEL_ID")
	token := os.Getenv("DISCORD_BOT_TOKEN")

	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Logger.Error(err.Error())
	}

	discord := Discord{
		bot:           bot,
		nickMemberMap: make(map[string]*discordgo.Member),
		messages:      NewMessages(512),
		message:       messageChan,
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
		logger.Fatal(fmt.Sprintf("failed to connect to discord: %v", err))
	}

	logger.Logger.Info("discord connection established")

	members, err := d.bot.GuildMembers(guildId, "", 1000)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error obtaining server members: %v", err.Error()))
	}

	for _, member := range members {
		if member == nil {
			logger.Logger.Warn("Skipping missing information for a user.")
			continue
		}

		d.nickMemberMap[member.User.Username] = member
		if member.Nick != "" {
			d.nickMemberMap[member.Nick] = member
		}
	}
}

func (d *Discord) Close() {
	err := d.bot.Close()
	if err != nil {
		logger.Logger.Error(err.Error())
	}
}

func (d *Discord) HandleMessage(msg message.Message) {
	switch msg.Type {
	case message.TypeDefault:
		d.sendMessage(fmt.Sprintf("<%s> %s", msg.Author, msg.Text))
	case message.TypeReply:
		d.replyToMessage(msg.Offset.Username, msg.Author, msg.Text, int32(msg.Offset.Offset))
	default:
		logger.Logger.Warn(fmt.Sprintf("unknown message type: %d", msg.Type))
	}
}

func (d *Discord) sendMessage(message string) {
	message = d.replaceUserMentions(message)
	msg, err := d.bot.ChannelMessageSend(channelId, message)
	if err != nil {
		logger.Logger.Error("error occurred during sending message", "err", err)
	}

	d.messages.push(MessageAuthor(message), Message{
		content:   MessageContent(message),
		messageId: msg.ID,
	})
}

func (d *Discord) onMessage(session *discordgo.Session, msg *discordgo.MessageCreate) {
	if session.State.User.ID != msg.Author.ID && msg.ChannelID == channelId && msg.Type == discordgo.MessageTypeDefault {
		author := msg.Author.Username
		content, err := msg.ContentWithMoreMentionsReplaced(session)
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("failed to parse discord message: %v", err))
		}
		content = replaceEmotes(content)

		d.message <- message.NewMessage(
			content,
			author,
			message.TargetIrc,
		)
		d.messages.push(author, Message{
			content:   content,
			messageId: msg.Message.ID,
		})
	} else if session.State.User.ID != msg.Author.ID && msg.ChannelID == channelId && msg.Type == discordgo.MessageTypeReply {
		d.onReply(session, msg)
	}
}

func (d *Discord) onReply(session *discordgo.Session, msg *discordgo.MessageCreate) {
	repliedMessage := msg.ReferencedMessage
	from := msg.Author.Username
	content := msg.Content
	to := ""

	if session.State.User.ID == repliedMessage.Author.ID {
		to = MessageAuthor(repliedMessage.Content)
	} else {
		to = repliedMessage.Author.Username
	}

	messageId := d.messages.find(to, repliedMessage.ID)

	d.message <- message.NewReplyMessage(
		content,
		from,
		message.TargetIrc,
		message.Offset{
			Offset:   int(messageId),
			Username: to,
		},
	)
	d.messages.push(from, Message{
		content:   content,
		messageId: msg.Message.ID,
	})
}

func (d *Discord) replyToMessage(username string, author string, content string, offset int32) {
	message := d.messages.findByOffset(username, uint32(offset))

	if message != nil {
		message, err := d.bot.ChannelMessageSendReply(channelId, fmt.Sprintf("<%s> %s", author, content), &discordgo.MessageReference{
			MessageID: message.messageId,
			ChannelID: channelId,
			GuildID:   guildId,
		})

		if err != nil {
			logger.Logger.Info("failed to reply to a message", "err", err)
			return
		}

		d.messages.push(username, Message{
			content:   content,
			messageId: message.ID,
		})
	}
}

func (d *Discord) onEdit(session *discordgo.Session, msg *discordgo.MessageUpdate) {
	author := msg.Author.Username
	content := msg.Content
	offset := d.messages.find(author, msg.ID)

	d.message <- message.NewEditMessage(
		content,
		author,
		message.TargetIrc,
		message.Offset{
			Offset:   int(offset),
			Username: author,
		},
	)
	d.messages.update(msg.ID, author, content)
}

func (d *Discord) onReactionAdd(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	emojiName := reaction.Emoji.Name
	from := reaction.Member.User.Username
	to := ""
	msg, err := session.ChannelMessage(channelId, reaction.MessageID)
	if err != nil {
		logger.Logger.Warn("failed to find message", "err", err)
		return
	}

	if session.State.User.ID == msg.Author.ID {
		to = MessageAuthor(msg.Content)
	} else {
		to = msg.Author.Username
	}

	messageId := d.messages.find(to, reaction.MessageID)

	d.message <- message.NewReactMessage(
		from,
		message.TargetIrc,
		message.Offset{
			Username: to,
			Offset:   int(messageId),
		}, message.Reaction{
			EmojiName: emojiName,
			Type:      message.ReactionTypeAdded,
		},
	)
}

func (d *Discord) onReactionRemove(session *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	emojiName := reaction.Emoji.Name
	user, err := session.User(reaction.UserID)
	if err != nil {
		logger.Logger.Warn("failed to find user", "err", err)
		return
	}
	from := user.Username
	to := ""
	msg, err := session.ChannelMessage(channelId, reaction.MessageID)
	if err != nil {
		logger.Logger.Warn("failed to get message", "err", err)
		return
	}

	if session.State.User.ID == msg.Author.ID {
		to = MessageAuthor(msg.Content)
	} else {
		to = msg.Author.Username
	}

	messageId := d.messages.find(to, reaction.MessageID)

	d.message <- message.NewReactMessage(
		from,
		message.TargetIrc,
		message.Offset{
			Username: to,
			Offset:   int(messageId),
		}, message.Reaction{
			EmojiName: emojiName,
			Type:      message.ReactionTypeRemoved,
		},
	)
}

func (d *Discord) deleteMessage(name string, offset int) {
	message := d.messages.findByOffset(name, uint32(offset))
	if message == nil {
		logger.Logger.Warn("failed to find message")
		return
	}

	d.bot.ChannelMessageDelete(channelId, message.messageId)
	d.messages.delete(name, message.messageId)
}

var (
	userMentionRE = regexp.MustCompile("@[^@\n ]{1,32}")
	emoteRE       = regexp.MustCompile(`<a?(:\w+:)\d+>`)
)

func (d *Discord) replaceUserMentions(content string) string {
	fn := func(match string) string {
		username := match[1:]
		member := d.nickMemberMap[username]

		if member == nil {
			return match
		}

		return strings.Replace(match, "@"+username, member.User.Mention(), 1)
	}

	return userMentionRE.ReplaceAllStringFunc(content, fn)
}

func replaceEmotes(text string) string {
	return emoteRE.ReplaceAllString(text, "$1")
}
