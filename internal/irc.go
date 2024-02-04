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
	} else {
		author := event.Nick
		content := event.Message()

		message := fmt.Sprintf("<%s> %s", author, content)

		i.bridge.discord.sendMessage(message)
	}
}

func (i *Irc) onReply(event *irc.Event) {
	message := event.Message()
	username := ""
	content := ""
	offset := 0
	j := 0

	for a := 7; a < len(message)-1; a++ {
		if message[a] == '~' {
			username = message[7:a]
			j = a + 1

			break
		}
	}

	for b := j; b < len(message)-1; b++ {
		if message[b] == ' ' {
			i, err := strconv.Atoi(message[j:b])
			if err != nil {
				fmt.Println("failed to convert string into a number: %w", err)
				return
			}

			offset = i
			content = message[b+1:]
		}
	}

	i.bridge.discord.replyToMessage(username, content, int32(offset))

}
