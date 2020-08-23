package reply

import (
	"math/rand"
	"strings"
)

type Replier interface {
	Reply(string) string
	CanReadMessage(payload *MessagePayload) bool
}

// Response defines a bot message response.
type Response struct {
	// Text is the reply message.
	Text string
	// Reply tells whether or not reply to the message.
	// See https://telegram.org/blog/replies-mentions.
	Reply bool
}

// MessagePayload defines a telegram message.
type MessagePayload struct {
	// Text is the content of the message.
	Text string
	// UserId contains the id of the message sender.
	UserId string
}

// RandInt is a function to return a random int. It's exposed to allow customization in tests.
var RandInt = rand.Int

// GetReply parses a message and responds accordingly.
// It can differ between a regular question and an indagation (see isIndagation method).
func GetReply(payload *MessagePayload) *Response {
	text := strings.ToLower(strings.TrimSpace(payload.Text))
	answer := getReplier(payload).Reply(text)

	if answer == "" && IsQuestion(text) {
		if IsIndagation(text) {
			answer = "sei la"
		} else {
			answer = "sim"
		}
	}

	return &Response{
		Text:  answer,
		Reply: RandInt()%2 == 0,
	}
}

var repliers []Replier

func RegisterReplier(r Replier) {
	repliers = append(repliers, r)
}

func getReplier(msg *MessagePayload) Replier {
	for _, replier := range repliers {
		if replier.CanReadMessage(msg) {
			return replier
		}
	}

	return &DefaultReplier{}
}