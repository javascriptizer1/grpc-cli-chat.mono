package converter

import (
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/domain"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/repository/dao"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToDaoMessageUser(domainUser *domain.MessageUser) *dao.MessageUser {
	return &dao.MessageUser{
		ID:   domainUser.ID,
		Name: domainUser.Name,
	}
}

func ToDomainMessageUser(daoUser *dao.MessageUser) *domain.MessageUser {
	return &domain.MessageUser{
		ID:   daoUser.ID,
		Name: daoUser.Name,
	}
}

func ToDaoMessage(domainMessage *domain.Message) (*dao.Message, error) {
	objectID, err := primitive.ObjectIDFromHex(domainMessage.ID)

	if err != nil {
		return nil, err
	}

	chatObjectID, err := primitive.ObjectIDFromHex(domainMessage.ChatID)

	if err != nil {
		return nil, err
	}

	return &dao.Message{
		ID:        objectID,
		Sender:    ToDaoMessageUser(domainMessage.Sender),
		ChatID:    chatObjectID,
		Text:      domainMessage.Text,
		CreatedAt: domainMessage.CreatedAt,
	}, nil
}

func ToDomainMessage(daoMessage *dao.Message) *domain.Message {
	return &domain.Message{
		ID:        daoMessage.ID.Hex(),
		Sender:    ToDomainMessageUser(daoMessage.Sender),
		ChatID:    daoMessage.ChatID.Hex(),
		Text:      daoMessage.Text,
		CreatedAt: daoMessage.CreatedAt,
	}
}

func ToDomainMessageList(daoMessages []*dao.Message) []*domain.Message {
	domainMessages := make([]*domain.Message, len(daoMessages))

	for i, message := range daoMessages {
		domainMessages[i] = ToDomainMessage(message)
	}

	return domainMessages
}
