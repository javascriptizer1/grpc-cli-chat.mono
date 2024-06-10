package userrepo

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.backend/internal/domain/user"

	"github.com/jmoiron/sqlx"
)

const (
	users = "users"
)

type UserRepository struct {
	db      *sqlx.DB
	builder sq.StatementBuilderType
}

func New(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *UserRepository) Create(_ context.Context, input *user.User) error {
	sql, args, err := r.builder.Insert(users).Columns(
		"id",
		"name",
		"email",
		"password_hash",
		"role",
		"created_at",
		"updated_at",
	).Values(
		input.Id,
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

	row := r.db.QueryRow(sql, args...)

	return row.Err()
}

func (r *UserRepository) OneById(_ context.Context, id uuid.UUID) (u *user.User, err error) {
	var rawUser User

	sql, args, err := r.builder.
		Select("id", "name", "email", "password_hash", "role", "created_at", "updated_at").
		From(users).
		Where(sq.Eq{"id": id.String()}).
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

func (r *UserRepository) OneByEmail(_ context.Context, email string) (u *user.User, err error) {
	var rawUser User

	sql, args, err := r.builder.
		Select("id", "name", "email", "password_hash", "role", "created_at", "updated_at").
		From(users).
		Where(sq.Eq{"email": email}).
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
