package pipeline

import (
	"time"
)

type Message struct {
	timestamp time.Time
	key       []byte
	value     []byte
}

func InitMessage() *Message {
	return new(Message)
}
