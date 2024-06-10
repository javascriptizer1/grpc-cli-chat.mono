package user

import (
	"time"

	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain"
)

type Role uint16

type UserRepo struct {
	Id           uint64    `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         Role      `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (u *UserRepo) ToDomain() *domain.User {
	return &domain.User{
		Id:           u.Id,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         domain.Role(u.Role),
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
