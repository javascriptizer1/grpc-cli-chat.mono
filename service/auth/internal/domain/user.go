package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.mono/pkg/helper/bcrypt"
	"github.com/javascriptizer1/grpc-cli-chat.mono/pkg/type/pagination"
)

type UserListFilter struct {
	pagination.Pagination
	UserIDs []uuid.UUID
}

type UserRole uint16

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUserWithID(id uuid.UUID, name string, email string, password string, role UserRole, createdAt time.Time, updatedAt time.Time) (*User, error) {
	user := &User{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  password,
		Role:      role,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return user, nil
}

func NewUser(name string, email string, password string, role UserRole) (*User, error) {
	now := time.Now().UTC()

	user := &User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  password,
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return user, nil
}

func (u *User) HashPassword() error {

	if u.Password == "" {
		return errors.New("password is empty")
	}

	p, err := bcrypt.Hash(u.Password)

	if err != nil {
		return err
	}

	u.Password = p

	return nil
}

func (u *User) CheckPassword(candidatePassword string) bool {
	return bcrypt.Check(candidatePassword, u.Password)
}
