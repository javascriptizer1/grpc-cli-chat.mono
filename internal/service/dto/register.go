package dto

import "github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain"

type RegisterInputDto struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            domain.UserRole
}
