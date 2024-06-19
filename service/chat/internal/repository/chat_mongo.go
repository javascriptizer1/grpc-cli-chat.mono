package repository

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.mono/pkg/type/pagination"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/chat/internal/converter"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/chat/internal/domain"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/chat/internal/repository/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	chatCollection = "chats"
)

type ChatRepository struct {
	db *mongo.Database
}

func NewChatRepository(db *mongo.Database) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Create(ctx context.Context, chat *domain.Chat) error {
	doc, err := converter.ToDaoChat(chat)

	if err != nil {
		return err
	}

	_, err = r.db.Collection(chatCollection).InsertOne(ctx, doc)

	return err
}

func (r *ChatRepository) OneByID(ctx context.Context, id string) (*domain.Chat, error) {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	var chat dao.Chat

	err = r.db.Collection(chatCollection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&chat)

	if err != nil {
		return nil, err
	}

	return converter.ToDomainChat(&chat), nil
}

func (r *ChatRepository) ContainUser(ctx context.Context, chatID string, userID string) bool {
	objectID, err := primitive.ObjectIDFromHex(chatID)

	if err != nil {
		return false
	}

	var chat dao.Chat

	err = r.db.Collection(chatCollection).FindOne(ctx, bson.M{"_id": objectID, "users.id": userID}).Decode(&chat)

	return err == nil
}

func (r *ChatRepository) List(ctx context.Context, userID string, p pagination.Pagination) ([]*domain.Chat, uint32, error) {
	filter := bson.M{"users.id": userID}

	options := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetLimit(int64(p.Limit())).
		SetSkip(int64(p.Offset()))

	cur, err := r.db.Collection(chatCollection).Find(ctx, filter, options)

	if err != nil {
		return nil, 0, err
	}

	defer func() {
		err = cur.Close(ctx)
	}()

	var chats []*dao.Chat

	if err = cur.All(ctx, &chats); err != nil {
		return nil, 0, err
	}

	total, err := r.db.Collection(chatCollection).CountDocuments(ctx, filter)

	if err != nil {
		return nil, 0, err
	}

	return converter.ToDomainChatList(chats), uint32(total), nil
}
