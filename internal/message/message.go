package message

type Target uint8

const (
	TargetDiscord Target = 0
	TargetIrc     Target = 1
)

type Type uint8

const (
	TypeDefault  Type = 0
	TypeReply    Type = 1
	TypeEdit     Type = 2
	TypeDelete   Type = 3
	TypeReaction Type = 4
)

type ReactionType uint8

const (
	ReactionTypeAdded   ReactionType = 0
	ReactionTypeRemoved ReactionType = 1
)

type Offset struct {
	Username string
	Offset   int
}

type Reaction struct {
	EmojiName string
	Type      ReactionType
}

type Message struct {
	Text     string
	Author   string
	Target   Target
	Type     Type
	Offset   *Offset
	Reaction *Reaction
}

func NewMessage(text string, author string, target Target) Message {
	return Message{
		Text:   text,
		Author: author,
		Target: target,
		Type:   TypeDefault,
	}
}

func NewReplyMessage(text string, Author string, target Target, offset Offset) Message {
	return Message{
		Text:   text,
		Author: Author,
		Target: target,
		Offset: &offset,
		Type:   TypeReply,
	}
}

func NewEditMessage(text string, author string, target Target, offset Offset) Message {
	return Message{
		Text:   text,
		Author: author,
		Target: target,
		Offset: &offset,
		Type:   TypeEdit,
	}
}

func NewDeleteMessage(text string, author string, target Target, offset Offset) Message {
	return Message{
		Text:   text,
		Author: author,
		Target: target,
		Offset: &offset,
		Type:   TypeDelete,
	}
}

func NewReactMessage(author string, target Target, offset Offset, reaction Reaction) Message {
	return Message{
		Author:   author,
		Target:   target,
		Offset:   &offset,
		Reaction: &reaction,
		Type:     TypeReaction,
	}
}
