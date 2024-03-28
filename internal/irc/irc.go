package irc

import (
	"fmt"
	"os"

	"strconv"
	"strings"

	"github.com/MilkeeyCat/deez_bridge/internal/logger"
	"github.com/MilkeeyCat/deez_bridge/internal/message"
	"github.com/thoj/go-ircevent"
)

type Irc struct {
	connection *irc.Connection
	message    chan message.Message
}

var channel string

func NewIrc(messageChan chan message.Message) Irc {
	nickname := os.Getenv("IRC_NICKNAME")
	channel = os.Getenv("IRC_CHANNEL")

	con := irc.IRC(nickname, nickname)
	con.UseTLS = false

	ircS := Irc{
		connection: con,
		message:    messageChan,
	}

	con.AddCallback("001", func(e *irc.Event) {
		con.Join(channel)
		logger.Logger.Info("irc connection established")
	})
	con.AddCallback("PRIVMSG", ircS.onMessage)

	return ircS
}

func (i *Irc) Open() {
	host := os.Getenv("IRC_SERVER_HOST")
	port := os.Getenv("IRC_SERVER_PORT")

	err := i.connection.Connect(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed connect to irc: %v", err))
	}
}

func (i *Irc) Close() {
	i.connection.Disconnect()
}

func (i *Irc) HandleMessage(msg message.Message) {
	switch msg.Type {
	case message.TypeDefault:
		i.sendMessage(fmt.Sprintf("<%s> %s", msg.Author, msg.Text))
		break
	case message.TypeReply:
		i.sendMessage(fmt.Sprintf("<%s %s~%d> %s", msg.Author, msg.Offset.Username, msg.Offset.Offset, msg.Text))
		break
	case message.TypeEdit:
		i.sendMessage(fmt.Sprintf("<%s~%d> %s", msg.Author, msg.Offset.Offset, msg.Text))
		break
	case message.TypeReaction:
		var content string

		if msg.Reaction.Type == message.ReactionTypeAdded {
			content = fmt.Sprintf("%s reacted with %s to %s~%d", msg.Author, msg.Reaction.EmojiName, msg.Offset.Username, msg.Offset.Offset)
		} else if msg.Reaction.Type == message.ReactionTypeRemoved {
			content = fmt.Sprintf("%s removed reaction %s from %s~%d", msg.Author, msg.Reaction.EmojiName, msg.Offset.Username, msg.Offset.Offset)
		} else {
			logger.Logger.Warn(fmt.Sprintf("unknown reaction type: %d", msg.Reaction.Type))
		}

		i.sendMessage(content)
		break
	default:
		logger.Logger.Warn(fmt.Sprintf("unknown message type: %d", msg.Type))
	}
}

func (i *Irc) sendMessage(message string) {
	i.connection.Privmsgf(channel, message)
}

func (i *Irc) onMessage(event *irc.Event) {
	//TODO: make it better
	if strings.HasPrefix(event.Message(), "!reply") {
		i.onReply(event)
	} else if strings.HasPrefix(event.Message(), "!del") {
		i.onDelete(event)
	} else {
		i.message <- message.NewMessage(
			event.Message(),
			event.Nick,
			message.TargetDiscord,
		)
	}
}

func (i *Irc) onReply(event *irc.Event) {
	msg := event.Message()
	vals := strings.SplitN(msg, " ", 3)
	if len(vals) != 3 {
		return
	}

	str := strings.Split(vals[1], "~")
	if len(str) != 2 {
		return
	}

	username := str[0]
	content := vals[2]
	offset, err := strconv.Atoi(str[1])
	if err != nil {
		logger.Logger.Info("failed to convert string", "str", str[1])
		return
	}

	i.message <- message.NewReplyMessage(
		content,
		event.Nick,
		message.TargetDiscord,
		message.Offset{
			Username: username,
			Offset:   offset,
		},
	)
}

func (i *Irc) onDelete(event *irc.Event) {
	msg := event.Message()
	vals := strings.Split(msg, " ")
	if len(vals) != 2 {
		return
	}

	offset, err := strconv.Atoi(vals[1])
	if err != nil {
		return
	}

	i.message <- message.NewDeleteMessage(
		msg,
		event.Nick,
		message.TargetDiscord,
		message.Offset{
			Username: event.Nick,
			Offset:   offset,
		},
	)
}
