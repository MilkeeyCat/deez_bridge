package bridge

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/thoj/go-ircevent"
)

type Irc struct {
	connection *irc.Connection
	bridge     *Bridge
}

var channel string

func NewIrcBridge(bridge *Bridge) *Irc {
	nickname := os.Getenv("IRC_NICKNAME")
	channel = os.Getenv("IRC_CHANNEL")

	con := irc.IRC(nickname, nickname)
	con.UseTLS = false

	ircS := &Irc{
		connection: con,
		bridge:     bridge,
	}

	con.AddCallback("001", func(e *irc.Event) {
		con.Join(channel)
	})
	con.AddCallback("PRIVMSG", ircS.onMessage)

	return ircS
}

func (i *Irc) Open() {
	host := os.Getenv("IRC_SERVER_HOST")
	port := os.Getenv("IRC_SERVER_PORT")

	err := i.connection.Connect(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		panic(err)
	}

	fmt.Println("irc connection established")
}

func (i *Irc) Close() {
	i.connection.Disconnect()
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
		author := event.Nick
		content := event.Message()

		message := fmt.Sprintf("<%s> %s", author, content)

		i.bridge.discord.sendMessage(message)
	}
}

func (i *Irc) onReply(event *irc.Event) {
	//TODO: add error handling
	message := event.Message()
	vals := strings.SplitN(message, " ", 3)
	str := strings.Split(vals[1], "~")
	username := str[0]
	content := vals[2]
	offset, err := strconv.Atoi(str[1])
	if err != nil {
		fmt.Printf("failed to convert string '%s' into int", str[1])
		return
	}

	i.bridge.discord.replyToMessage(username, event.Nick, content, int32(offset))
}

func (i *Irc) onDelete(event *irc.Event) {
	//TOOD: add error handling
	message := event.Message()
	vals := strings.Split(message, " ")
	offset, err := strconv.Atoi(vals[1])
	if err != nil {
		fmt.Printf("failed to convert string '%s' into int", vals[1])
		return
	}

	i.bridge.discord.deleteMessage(event.Nick, offset)
}
