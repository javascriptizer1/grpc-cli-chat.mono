package dao

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageUser struct {
	ID   string `bson:"id"`
	Name string `bson:"name"`
}

type Message struct {
	ID        primitive.ObjectID `bson:"_id"`
	Sender    *MessageUser       `bson:"sender"`
	ChatID    primitive.ObjectID `bson:"chatId"`
	Text      string             `bson:"text"`
	CreatedAt time.Time          `bson:"createdAt"`
}

func NewMessage(objectID, chatObjectID primitive.ObjectID, sender *MessageUser, text string, createdAt time.Time) *Message {
	return &Message{
		ID:        objectID,
		Sender:    sender,
		ChatID:    chatObjectID,
		Text:      text,
		CreatedAt: createdAt,
	}
}

func NewMessageUser(id, name string) *MessageUser {
	return &MessageUser{
		ID:   id,
		Name: name,
	}
}
