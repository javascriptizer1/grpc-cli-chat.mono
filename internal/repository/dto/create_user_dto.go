package dto

import repository "github.com/javascriptizer1/grpc-cli-chat.backend/internal/repository/model"

type CreateUserDto struct {
	Name     string
	Email    string
	Password string
	Role     repository.Role
}
