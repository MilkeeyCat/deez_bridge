package discord

import "strings"

type Message struct {
	messageId string
	content   string
}

type MessagesMap map[string][]Message

type Messages struct {
	messages MessagesMap
	maxMsgs  int
}

func (m *Messages) push(author string, msg Message) {
	userMessages := m.messages[author]

	userMessages = append(userMessages, msg)
	m.messages[author] = userMessages
}

func (m *Messages) find(user, messageId string) int32 {
	msgs := m.messages[user]

	for i := range msgs {
		if msgs[i].messageId == messageId {
			return int32(len(msgs) - i)
		}
	}

	return -1
}

func (m *Messages) update(messageId, user, msg string) {
	msgs := m.messages[user]

	for i := range msgs {
		if msgs[i].messageId == messageId {
			m.messages[user][i].content = msg
		}
	}
}

func (m *Messages) delete(user, messageId string) {
	msgs := m.messages[user]

	for i := range msgs {
		if msgs[i].messageId == messageId {
			m.messages[user] = append(msgs[:i], msgs[i+1:]...)
		}
	}
}

func (m *Messages) findByOffset(user string, offset uint32) *Message {
	msgs := m.messages[user]

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
