package service

import (
	"context"

	client "github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/client/grpc"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/domain"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/chat/internal/repository"
)

type ChatService struct {
	chatRepo   repository.ChatRepository
	userClient client.UserClient
}

func NewChatService(chatRepo repository.ChatRepository, userClient client.UserClient) *ChatService {
	return &ChatService{
		chatRepo:   chatRepo,
		userClient: userClient,
	}
}

func (s *ChatService) Create(ctx context.Context, userID string, userIDs []string) (string, error) {
	users := make([]domain.ChatUser, 5)
	c := domain.NewChat(users)

	err := s.chatRepo.Create(ctx, c)

	return c.ID, err
}

func (s *ChatService) List(ctx context.Context, userID string) ([]*domain.Chat, error) {
	return s.chatRepo.List(ctx, userID)
}
