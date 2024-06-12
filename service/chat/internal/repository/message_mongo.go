package repository

import (
	"context"

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
	doc, err := dao.ToDaoMessage(*message)

	if err != nil {
		return err
	}

	_, err = r.db.Collection(messageCollection).InsertOne(ctx, doc)

	if err != nil {
		return err
	}

	return nil
}

func (r *MessageRepository) List(ctx context.Context, userID string) ([]*domain.Message, int, error) {

	var res []*dao.Message

	filter := bson.M{"users.id": bson.M{"in": userID}}

	cur, err := r.db.
		Collection(messageCollection).
		Find(ctx,
			filter,
			&options.FindOptions{Sort: bson.M{"createdAt": -1}},
		)

	if err != nil {
		return []*domain.Message{}, 0, err
	}

	total, err := r.db.Collection(messageCollection).CountDocuments(ctx, filter)

	if err != nil {
		return []*domain.Message{}, 0, err
	}

	err = cur.Decode(&res)

	if err != nil {
		return []*domain.Message{}, 0, err
	}

	return dao.ToDomainMessageList(res), int(total), nil
}
