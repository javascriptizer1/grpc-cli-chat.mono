package repository

import (
	"context"

	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/type/pagination"
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
	doc, err := dao.ToDaoChat(chat.ID, chat.Name, chat.Users, chat.CreatedAt)

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

	var c *dao.Chat

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	res := r.db.Collection(chatCollection).FindOne(ctx,
		bson.M{"_id": objectID},
	)

	err = res.Decode(&c)

	return dao.ToDomainChat(*c), err

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

func (r *ChatRepository) List(ctx context.Context, userID string, p pagination.Pagination) ([]*domain.Chat, uint32, error) {

	var res []dao.Chat

	filter := bson.M{"users.id": bson.M{"$in": []string{userID}}}

	cur, err := r.db.
		Collection(chatCollection).
		Find(ctx,
			filter,
			options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}),
			options.Find().SetLimit(int64(p.Limit())),
			options.Find().SetSkip(int64(p.Offset())),
		)

	if err != nil {
		return []*domain.Chat{}, 0, err
	}

	total, err := r.db.Collection(chatCollection).CountDocuments(ctx, filter)

	if err != nil {
		return []*domain.Chat{}, 0, err
	}

	err = cur.All(ctx, &res)

	if err != nil {
		return []*domain.Chat{}, 0, err
	}

	return dao.ToDomainChatList(res), uint32(total), nil
}
