package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatUser struct {
	ID          string
	ConnectedAt time.Time
}

type Chat struct {
	ID        string
	Name      string
	Users     []ChatUser
	CreatedAt time.Time
}

func NewChat(name string, users []ChatUser) *Chat {
	c := &Chat{
		ID:        primitive.NewObjectID().Hex(),
		Name:      name,
		Users:     users,
		CreatedAt: time.Now().UTC(),
	}

	return c
}

func NewChatUser(id string) *ChatUser {
	cu := &ChatUser{
		ID:          id,
		ConnectedAt: time.Now().UTC(),
	}

	return cu
}
