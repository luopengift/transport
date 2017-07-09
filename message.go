package transport

import (
	"time"
)

type Message struct {
	timestamp time.Time
	key       []byte
	value     []byte
}
