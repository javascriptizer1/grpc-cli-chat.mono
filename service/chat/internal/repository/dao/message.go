package dao

import (
	"time"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageUser struct {
	ID   string `bson:"id"`
	Name string `bson:"name"`
}

type Message struct {
	ID        primitive.ObjectID `bson:"_id"`
	Sender    MessageUser        `bson:"sender"`
	ChatID    primitive.ObjectID `bson:"chatId"`
	Text      string             `bson:"text"`
	CreatedAt time.Time          `bson:"createdAt"`
}

func ToDaoMessage(message domain.Message) (*Message, error) {
	objectID, err := primitive.ObjectIDFromHex(message.ID)

	if err != nil {
		return nil, err
	}

	chatObjectID, err := primitive.ObjectIDFromHex(message.ChatID)

	if err != nil {
		return nil, err
	}

	dm := &Message{
		ID:        objectID,
		Sender:    *ToDaoMessageSender(message.Sender),
		ChatID:    chatObjectID,
		Text:      message.Text,
		CreatedAt: message.CreatedAt,
	}

	return dm, err
}

func ToDaoMessageSender(user domain.MessageUser) *MessageUser {
	dmu := &MessageUser{
		ID:   user.ID,
		Name: user.Name,
	}

	return dmu
}

func ToDomainMessage(m Message) *domain.Message {

	dm := &domain.Message{
		ID:        m.ID.Hex(),
		Sender:    *ToDomainMessageSender(m.Sender),
		ChatID:    m.ChatID.Hex(),
		Text:      m.Text,
		CreatedAt: m.CreatedAt,
	}

	return dm
}

func ToDomainMessageList(dao []*Message) []*domain.Message {

	var result = make([]*domain.Message, len(dao))

	for i, v := range dao {
		cu := ToDomainMessage(*v)

		result[i] = cu
	}

	return result
}

func ToDomainMessageSender(user MessageUser) *domain.MessageUser {
	dmu := &domain.MessageUser{
		ID:   user.ID,
		Name: user.Name,
	}

	return dmu
}
