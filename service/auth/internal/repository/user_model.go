package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/auth/internal/domain"
)

type userRole uint16

type User struct {
	ID           string    `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         userRole  `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (u *User) ToDomain() *domain.User {
	user, _ := domain.NewUserWithID(
		uuid.MustParse(u.ID),
		u.Name,
		u.Email,
		u.PasswordHash,
		domain.UserRole(u.Role),
		u.CreatedAt,
		u.UpdatedAt,
	)

	return user
}

func ListToDomain(users []*User) []*domain.User {
	du := make([]*domain.User, len(users))

	for i, v := range users {
		du[i] = v.ToDomain()
	}

	return du
}
