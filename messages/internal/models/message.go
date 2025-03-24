package models

import "time"

type Message struct {
	ChatID    string
	MessageID string
	SenderID  string
	Content   string
	Timestamp time.Time
}
