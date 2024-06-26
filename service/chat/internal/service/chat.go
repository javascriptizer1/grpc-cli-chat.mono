package service

import (
	"context"
	"errors"

	"github.com/javascriptizer1/grpc-cli-chat.mono/pkg/type/pagination"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/chat/internal/domain"
)

type ChatService struct {
	chatRepo    ChatRepository
	messageRepo MessageRepository
	userClient  UserClient
}

func NewChatService(chatRepo ChatRepository, messageRepo MessageRepository, userClient UserClient) *ChatService {
	return &ChatService{
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
		userClient:  userClient,
	}
}

func (s *ChatService) Create(ctx context.Context, name string, userIDs []string) (string, error) {
	users := make([]*domain.ChatUser, len(userIDs))

	for i, id := range userIDs {
		users[i] = domain.NewChatUser(id)
	}

	c := domain.NewChat(name, users)

	err := s.chatRepo.Create(ctx, c)

	return c.ID, err
}

func (s *ChatService) OneByID(ctx context.Context, id string) (*domain.Chat, error) {
	return s.chatRepo.OneByID(ctx, id)
}

func (s *ChatService) List(ctx context.Context, userID string, p pagination.Pagination) ([]*domain.Chat, uint32, error) {
	return s.chatRepo.List(ctx, userID, p)
}

func (s *ChatService) CreateMessage(ctx context.Context, text string, chatID string, userInfo domain.UserInfo) (*domain.Message, error) {
	access := s.chatRepo.ContainUser(ctx, chatID, userInfo.ID)

	if !access {
		return nil, errors.New("chat not found")
	}

	m := domain.NewMessage(chatID, text, &domain.MessageUser{
		ID:   userInfo.ID,
		Name: userInfo.Name,
	})

	err := s.messageRepo.Create(ctx, m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *ChatService) ListMessage(ctx context.Context, chatID string, userID string) ([]*domain.Message, int, error) {

	ok := s.chatRepo.ContainUser(ctx, chatID, userID)

	if !ok {
		return []*domain.Message{}, 0, errors.New("forbidden to read another chat")
	}

	return s.messageRepo.List(ctx, chatID)
}
