package domain

import "time"

type Role uint16

type User struct {
	Id           uint64
	Name         string
	Email        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
