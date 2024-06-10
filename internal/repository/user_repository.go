package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/repository/dto"
	user "github.com/javascriptizer1/grpc-cli-chat.backend/internal/repository/model"
	"github.com/jmoiron/sqlx"
)

const (
	users = "users"
)

type UserRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(_ context.Context, input dto.CreateUserDto) (uint64, error) {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Insert(users).Columns(
		"name",
		"email",
		"password_hash",
		"role",
	).Values(
		input.Name,
		input.Email,
		input.Password,
		input.Role,
	).Suffix("RETURNING \"id\"").ToSql()

	if err != nil {
		return 0, err
	}

	row := r.db.QueryRow(sql, args...)

	var id uint64

	if err = row.Scan(&id); err != nil {
		return 0, err
	}

	return id, err
}

func (r *UserRepository) One(_ context.Context, id uint64) (u *domain.User, err error) {
	var rawUser user.UserRepo

	sql, args, err := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("id", "name", "email", "password_hash", "role", "created_at", "updated_at").
		From(users).
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return u, err
	}

	err = r.db.Get(&rawUser, sql, args...)

	if err != nil {
		return u, err
	}

	return rawUser.ToDomain(), nil
}
