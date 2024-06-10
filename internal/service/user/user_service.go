package usersvc

import (
	"context"

	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain/user"
)

func (s *UserService) OneById(ctx context.Context, id uuid.UUID) (*user.User, error) {
	return s.userRepo.OneById(ctx, id)
}
