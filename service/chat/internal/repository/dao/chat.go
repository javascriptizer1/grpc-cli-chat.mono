package dao

import (
	"time"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatUser struct {
	ID          string    `bson:"id"`
	Name        string    `bson:"name"`
	ConnectedAt time.Time `bson:"connectedAt"`
}

type Chat struct {
	ID        primitive.ObjectID `bson:"_id"`
	Users     []ChatUser         `bson:"users"`
	CreatedAt time.Time          `bson:"createdAt"`
}

func ToDaoChat(id string, users []domain.ChatUser, createdAt time.Time) (*Chat, error) {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	dc := &Chat{
		ID:        objectID,
		Users:     ToDaoChatUserList(users),
		CreatedAt: createdAt,
	}

	return dc, err
}

func ToDomainChat(dao Chat) *domain.Chat {

	dc := &domain.Chat{
		ID:        dao.ID.Hex(),
		Users:     ToDomainChatUsers(dao.Users),
		CreatedAt: dao.CreatedAt,
	}

	return dc
}

func ToDomainChatList(dao []*Chat) []*domain.Chat {
	var result = make([]*domain.Chat, len(dao))

	for i, v := range dao {
		cu := ToDomainChat(*v)

		result[i] = cu
	}

	return result
}

func ToDaoChatUser(domain *domain.ChatUser) *ChatUser {
	dcu := &ChatUser{
		ID:          domain.ID,
		Name:        domain.Name,
		ConnectedAt: domain.ConnectedAt,
	}

	return dcu
}

func ToDaoChatUserList(domain []domain.ChatUser) []ChatUser {
	var result = make([]ChatUser, len(domain))

	for i, v := range domain {
		cu := ToDaoChatUser(&v)

		result[i] = *cu
	}

	return result
}

func ToDomainChatUser(dao *ChatUser) *domain.ChatUser {
	dcu := &domain.ChatUser{
		ID:          dao.ID,
		Name:        dao.Name,
		ConnectedAt: dao.ConnectedAt,
	}

	return dcu
}

func ToDomainChatUsers(dao []ChatUser) []domain.ChatUser {
	var result = make([]domain.ChatUser, len(dao))

	for i, v := range dao {
		cu := ToDomainChatUser(&v)

		result[i] = *cu
	}

	return result
}
