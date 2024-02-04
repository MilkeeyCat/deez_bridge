package bridge

type Bridge struct {
	irc      *Irc
	discord  *Discord
	messages *MessagesMap
}

func NewBridge() Bridge {
    messagesMap := make(MessagesMap)
	bridge := Bridge{
		irc:      nil,
		discord:  nil,
		messages: &messagesMap,
	}

	bridge.irc = NewIrcBridge(&bridge)
	bridge.discord = NewDiscordBridge(&bridge)

	return bridge
}

func (b *Bridge) Open() {
	b.irc.Open()
	b.discord.Open()
}

func (b *Bridge) Close() {
	b.irc.Close()
	b.discord.Close()
}
