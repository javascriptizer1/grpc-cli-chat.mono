package dto

import "github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain/user"

type RegisterInputDto struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            user.Role
}
