package dao

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatUser struct {
	ID          string    `bson:"id"`
	ConnectedAt time.Time `bson:"connectedAt"`
}

type Chat struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Users     []*ChatUser        `bson:"users"`
	CreatedAt time.Time          `bson:"createdAt"`
}

func NewChat(objectID primitive.ObjectID, name string, users []*ChatUser, createdAt time.Time) *Chat {
	return &Chat{
		ID:        objectID,
		Name:      name,
		Users:     users,
		CreatedAt: createdAt,
	}
}

func NewChatUser(id string, connectedAt time.Time) *ChatUser {
	return &ChatUser{
		ID:          id,
		ConnectedAt: connectedAt,
	}
}
