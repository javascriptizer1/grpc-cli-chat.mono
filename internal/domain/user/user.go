package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/helper/bcrypt"
)

type Role uint16

type User struct {
	Id        uuid.UUID
	Name      string
	Email     string
	Password  string
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewWithID(id uuid.UUID, name string, email string, password string, role Role, createdAt time.Time, updatedAt time.Time) (*User, error) {
	user := &User{
		Id:        id,
		Name:      name,
		Email:     email,
		Password:  password,
		Role:      role,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return user, nil
}

func New(name string, email string, password string, role Role) (*User, error) {
	now := time.Now().UTC()

	user := &User{
		Id:        uuid.New(),
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
