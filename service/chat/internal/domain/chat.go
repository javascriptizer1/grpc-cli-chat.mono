package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatUser struct {
	ID          string
	Name        string
	ConnectedAt time.Time
}

type Chat struct {
	ID        string
	Users     []ChatUser
	CreatedAt time.Time
}

func NewChat(users []ChatUser) *Chat {
	c := &Chat{
		ID:        primitive.NewObjectID().String(),
		Users:     users,
		CreatedAt: time.Now().UTC(),
	}

	return c
}

func NewChatWithId(id string, users []ChatUser, createdAt time.Time) *Chat {
	c := &Chat{
		ID:        id,
		Users:     users,
		CreatedAt: createdAt,
	}

	return c
}
