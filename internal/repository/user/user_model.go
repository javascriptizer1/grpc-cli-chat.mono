package userrepo

import (
	"time"

	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain/user"
)

type Role uint16

type User struct {
	Id           string    `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         Role      `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (u *User) ToDomain() *user.User {
	return &user.User{
		Id:        uuid.MustParse(u.Id),
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.PasswordHash,
		Role:      user.Role(u.Role),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
