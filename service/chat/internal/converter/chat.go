package converter

import (
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/domain"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/repository/dao"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToDaoChatUser(domainUser *domain.ChatUser) *dao.ChatUser {
	return &dao.ChatUser{
		ID:          domainUser.ID,
		ConnectedAt: domainUser.ConnectedAt,
	}
}

func ToDomainChatUser(daoUser *dao.ChatUser) *domain.ChatUser {
	return &domain.ChatUser{
		ID:          daoUser.ID,
		ConnectedAt: daoUser.ConnectedAt,
	}
}

func ToDaoChatUserList(domainUsers []*domain.ChatUser) []*dao.ChatUser {
	daoUsers := make([]*dao.ChatUser, len(domainUsers))

	for i, user := range domainUsers {
		daoUsers[i] = ToDaoChatUser(user)
	}

	return daoUsers
}

func ToDomainChatUserList(daoUsers []*dao.ChatUser) []*domain.ChatUser {
	domainUsers := make([]*domain.ChatUser, len(daoUsers))

	for i, user := range daoUsers {
		domainUsers[i] = ToDomainChatUser(user)
	}

	return domainUsers
}

func ToDaoChat(domainChat *domain.Chat) (*dao.Chat, error) {
	objectID, err := primitive.ObjectIDFromHex(domainChat.ID)

	if err != nil {
		return nil, err
	}

	return &dao.Chat{
		ID:        objectID,
		Name:      domainChat.Name,
		Users:     ToDaoChatUserList(domainChat.Users),
		CreatedAt: domainChat.CreatedAt,
	}, nil
}

func ToDomainChat(daoChat *dao.Chat) *domain.Chat {
	return &domain.Chat{
		ID:        daoChat.ID.Hex(),
		Name:      daoChat.Name,
		Users:     ToDomainChatUserList(daoChat.Users),
		CreatedAt: daoChat.CreatedAt,
	}
}

func ToDomainChatList(daoChats []*dao.Chat) []*domain.Chat {
	domainChats := make([]*domain.Chat, len(daoChats))

	for i, chat := range daoChats {
		domainChats[i] = ToDomainChat(chat)
	}

	return domainChats
}
