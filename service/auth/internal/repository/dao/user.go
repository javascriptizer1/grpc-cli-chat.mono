package dao

import (
	"time"
)

type UserRole uint16

type User struct {
	ID           string    `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         UserRole  `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func NewUser(id, name, email, passwordHash string, role UserRole, createdAt, updatedAt time.Time) *User {
	return &User{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}
