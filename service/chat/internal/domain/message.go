package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageUser struct {
	ID   string
	Name string
}

type Message struct {
	ID        string
	Sender    *MessageUser
	ChatID    string
	Text      string
	CreatedAt time.Time
}

func NewMessage(chatID string, text string, sender *MessageUser) *Message {
	return &Message{
		ID:     primitive.NewObjectID().Hex(),
		ChatID: chatID,
		Text:   text,
		Sender: &MessageUser{
			ID:   sender.ID,
			Name: sender.Name,
		},
		CreatedAt: time.Now(),
	}
}
