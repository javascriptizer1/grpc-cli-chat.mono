package repository

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/converter"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/domain"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/repository/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	messageCollection = "messages"
)

type MessageRepository struct {
	db *mongo.Database
}

func NewMessageRepository(db *mongo.Database) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, message *domain.Message) error {
	doc, err := converter.ToDaoMessage(message)

	if err != nil {
		return err
	}

	_, err = r.db.Collection(messageCollection).InsertOne(ctx, doc)

	return err
}

func (r *MessageRepository) List(ctx context.Context, userID string) ([]*domain.Message, int, error) {
	filter := bson.M{"users.id": userID}

	options := options.Find().SetSort(bson.M{"createdAt": -1})

	cur, err := r.db.Collection(messageCollection).Find(ctx, filter, options)

	if err != nil {
		return nil, 0, err
	}

	defer func() {
		err = cur.Close(ctx)
	}()

	var messages []*dao.Message

	if err = cur.All(ctx, &messages); err != nil {
		return nil, 0, err
	}

	total, err := r.db.Collection(messageCollection).CountDocuments(ctx, filter)

	if err != nil {
		return nil, 0, err
	}

	return converter.ToDomainMessageList(messages), int(total), nil
}
