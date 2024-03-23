package bridge

import (
	"github.com/MilkeeyCat/deez_bridge/internal/discord"
	"github.com/MilkeeyCat/deez_bridge/internal/irc"
	"github.com/MilkeeyCat/deez_bridge/internal/message"
)

type Bridge struct {
	irc     irc.Irc
	discord discord.Discord
	message chan message.Message
}

func NewBridge() Bridge {
	messageChan := make(chan message.Message)

	bridge := Bridge{
		irc:     irc.NewIrc(messageChan),
		discord: discord.NewDiscord(messageChan),
		message: messageChan,
	}

	return bridge
}

func (b *Bridge) Open() {
	b.irc.Open()
	b.discord.Open()

	go b.Run()
}

func (b *Bridge) Run() {
	for msg := range b.message {
		if msg.Target == message.TargetDiscord {
			b.discord.HandleMessage(msg)
		} else if msg.Target == message.TargetIrc {
			b.irc.HandleMessage(msg)
		}
	}
}

func (b *Bridge) Close() {
	b.irc.Close()
	b.discord.Close()
}
