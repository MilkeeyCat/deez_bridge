package message

type target uint8

const (
	TargetDiscord target = iota
	TargetIrc
)

type _type uint8

const (
	TypeDefault _type = iota
	TypeReply
	TypeEdit
	TypeDelete
	TypeReaction
)

type reactionType uint8

const (
	ReactionTypeAdded reactionType = iota
	ReactionTypeRemoved
)

type Offset struct {
	Username string
	Offset   int
}

type Reaction struct {
	EmojiName string
	Type      reactionType
}

type Message struct {
	Text     string
	Author   string
	Target   target
	Type     _type
	Offset   *Offset
	Reaction *Reaction
}

func NewMessage(text string, author string, target target) Message {
	return Message{
		Text:   text,
		Author: author,
		Target: target,
		Type:   TypeDefault,
	}
}

func NewReplyMessage(text string, Author string, target target, offset Offset) Message {
	return Message{
		Text:   text,
		Author: Author,
		Target: target,
		Offset: &offset,
		Type:   TypeReply,
	}
}

func NewEditMessage(text string, author string, target target, offset Offset) Message {
	return Message{
		Text:   text,
		Author: author,
		Target: target,
		Offset: &offset,
		Type:   TypeEdit,
	}
}

func NewDeleteMessage(text string, author string, target target, offset Offset) Message {
	return Message{
		Text:   text,
		Author: author,
		Target: target,
		Offset: &offset,
		Type:   TypeDelete,
	}
}

func NewReactMessage(author string, target target, offset Offset, reaction Reaction) Message {
	return Message{
		Author:   author,
		Target:   target,
		Offset:   &offset,
		Reaction: &reaction,
		Type:     TypeReaction,
	}
}
