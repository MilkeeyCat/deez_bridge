package bridge

import (
	"sync"

	"github.com/MilkeeyCat/deez_bridge/internal/discord"
	"github.com/MilkeeyCat/deez_bridge/internal/irc"
)

type Bridge struct {
	irc     irc.BridgeIrc
	discord discord.BridgeDiscord
}

func NewBridge() Bridge {
	return Bridge{
		irc:     irc.NewIrcBridge(),
		discord: discord.NewDiscordBridge(),
	}
}

func (bridge Bridge) Run(wg *sync.WaitGroup) {
	go bridge.discord.Run()
    wg.Add(1)

	go bridge.irc.Run()
    wg.Add(1)
}
