package repository

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/domain"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/repository/dao"
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
	doc, err := dao.ToDaoChat(chat.ID, chat.Users, chat.CreatedAt)

	if err != nil {
		return err
	}

	_, err = r.db.Collection(chatCollection).InsertOne(ctx, doc)

	if err != nil {
		return err
	}

	return nil
}

func (r *ChatRepository) OneByID(ctx context.Context, id string) (*domain.Chat, error) {

	var c *domain.Chat

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c, err
	}

	ur := r.db.Collection(chatCollection).FindOne(ctx,
		bson.M{"_id": objectID},
	)

	err = ur.Decode(&c)

	return c, err

}

func (r *ChatRepository) ContainUser(ctx context.Context, chatID string, userID string) bool {
	var c *dao.Chat
	objectID, err := primitive.ObjectIDFromHex(chatID)

	if err != nil {
		return false
	}

	res := r.db.
		Collection(chatCollection).
		FindOne(ctx,
			bson.M{"_id": objectID, "users.id": bson.M{"$in": []string{userID}}},
		)

	err = res.Decode(&c)

	return c != nil && err == nil
}

func (r *ChatRepository) List(ctx context.Context, userID string) ([]*domain.Chat, error) {

	var res []*dao.Chat

	cur, err := r.db.
		Collection(chatCollection).
		Find(ctx,
			bson.M{"users.id": bson.M{"$in": []string{userID}}},
			&options.FindOptions{Sort: bson.M{"createdAt": -1}},
		)

	if err != nil {
		return []*domain.Chat{}, err
	}

	err = cur.Decode(&res)

	if err != nil {
		return []*domain.Chat{}, err
	}

	return dao.ToDomainChatList(res), nil
}
