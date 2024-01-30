package bridge

import "sync"

type Bridge struct {
	irc     *Irc
	discord *Discord
}

func NewBridge() Bridge {
	irc := NewIrcBridge()
	discord := NewDiscordBridge()

	irc.setDiscord(&discord)
	discord.setIrc(&irc)

	irc.Setup()
	discord.Setup()

	return Bridge{
		irc:     &irc,
		discord: &discord,
	}
}

func (b *Bridge) Run() {
	var wg sync.WaitGroup

	go b.discord.Run()
	wg.Add(1)

	go b.irc.Run()
	wg.Add(1)

	// where am i even supposed to call wg.Done()? xd
	wg.Wait()
}
