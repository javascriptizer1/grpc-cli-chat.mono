package repository

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/auth/internal/domain"

	"github.com/jmoiron/sqlx"
)

const (
	users = "users"
)

type UserRepository struct {
	db      *sqlx.DB
	builder sq.StatementBuilderType
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *UserRepository) Create(_ context.Context, input *domain.User) error {
	query, args, err := r.builder.Insert(users).Columns(
		"id",
		"name",
		"email",
		"password_hash",
		"role",
		"created_at",
		"updated_at",
	).Values(
		input.ID,
		input.Name,
		input.Email,
		input.Password,
		input.Role,
		input.CreatedAt,
		input.UpdatedAt,
	).ToSql()

	if err != nil {
		return err
	}

	row := r.db.QueryRow(query, args...)

	return row.Err()
}

func (r *UserRepository) OneByID(_ context.Context, id uuid.UUID) (u *domain.User, err error) {
	var rawUser User

	query, args, err := r.builder.
		Select("id", "name", "email", "password_hash", "role", "created_at", "updated_at").
		From(users).
		Where(sq.Eq{"id": id.String()}).
		ToSql()

	if err != nil {
		return u, err
	}

	err = r.db.Get(&rawUser, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, errors.New("user not found")
		}

		return u, err
	}

	return rawUser.ToDomain(), nil
}

func (r *UserRepository) OneByEmail(_ context.Context, email string) (u *domain.User, err error) {
	var rawUser User

	query, args, err := r.builder.
		Select("id", "name", "email", "password_hash", "role", "created_at", "updated_at").
		From(users).
		Where(sq.Eq{"email": email}).
		ToSql()

	if err != nil {
		return u, err
	}

	err = r.db.Get(&rawUser, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, errors.New("user not found")
		}

		return u, err
	}

	return rawUser.ToDomain(), nil
}
