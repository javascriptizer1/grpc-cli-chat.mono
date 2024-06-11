package domain

import "time"

type MessageUser struct {
	ID   string
	Name string
}

type Message struct {
	ID        string
	Sender    MessageUser
	ChatID    string
	Text      string
	CreatedAt time.Time
}
