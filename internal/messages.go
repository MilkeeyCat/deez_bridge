package bridge

import "strings"

type Message struct {
	messageId string
	content   string
}

type MessagesMap map[string][]Message

func (mm *MessagesMap) push(author string, msg Message) {
	userMessages := (*mm)[author]

	userMessages = append(userMessages, msg)
	(*mm)[author] = userMessages
}

func (mm *MessagesMap) find(user, messageId string) int32 {
	msgs := (*mm)[user]

	for i := range msgs {
		if msgs[i].messageId == messageId {
			return int32(len(msgs) - i)
		}
	}

	return -1
}

func (mm *MessagesMap) findByOffset(user string, offset uint32) *Message {
	msgs := (*mm)[user]

	if len(msgs) >= int(offset) {
		return &msgs[len(msgs)-int(offset)]
	}

	return nil
}

func MessageContent(msg string) string {
	cb := strings.Index(msg, ">")

	return msg[cb+2:]
}

func MessageAuthor(msg string) string {
	caretIndex := strings.Index(msg, "^")

	if caretIndex == -1 {
		cb := strings.Index(msg, ">")

		return msg[1:cb]
	}

	for i := range msg {
		if msg[i] == ' ' {
			return msg[1:i]
		}
	}

	return ""
}
