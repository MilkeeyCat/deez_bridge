package bridge

import "strings"

type MessagesMap map[string][]string

func (mm *MessagesMap) push(author, msg string) {
	userMessages := (*mm)[author]

	userMessages = append(userMessages, msg)
	(*mm)[author] = userMessages
}

func (mm *MessagesMap) find(user, msg string) int32 {
	msgs := (*mm)[user]

	for i := range msgs {
		if msgs[i] == msg {
			return int32(len(msgs) - i)
		}
	}

	return -1
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
