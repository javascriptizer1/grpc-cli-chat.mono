package repository

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/auth/internal/converter"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/auth/internal/domain"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/auth/internal/repository/dao"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"
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

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query, args, err := r.builder.Insert(usersTable).Columns(
		"id", "name", "email", "password_hash", "role", "created_at", "updated_at",
	).Values(
		user.ID, user.Name, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt,
	).ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *UserRepository) OneByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var rawUser dao.User

	query, args, err := r.builder.
		Select("id", "name", "email", "password_hash", "role", "created_at", "updated_at").
		From(usersTable).
		Where(sq.Eq{"id": id.String()}).
		ToSql()

	if err != nil {
		return nil, err
	}

	err = r.db.GetContext(ctx, &rawUser, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	return converter.ToDomainUser(&rawUser), nil
}

func (r *UserRepository) OneByEmail(ctx context.Context, email string) (*domain.User, error) {
	var rawUser dao.User

	query, args, err := r.builder.
		Select("id", "name", "email", "password_hash", "role", "created_at", "updated_at").
		From(usersTable).
		Where(sq.Eq{"email": email}).
		ToSql()

	if err != nil {
		return nil, err
	}

	err = r.db.GetContext(ctx, &rawUser, query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	return converter.ToDomainUser(&rawUser), nil
}

func (r *UserRepository) List(ctx context.Context, filter *domain.UserListFilter) ([]*domain.User, error) {
	var rawUsers []*dao.User

	builder := r.builder.
		Select("id", "name", "email", "password_hash", "role", "created_at", "updated_at").
		From(usersTable).
		Limit(uint64(filter.Limit())).
		Offset(uint64(filter.Offset())).
		OrderBy("updated_at desc")

	if len(filter.UserIDs) != 0 {
		uIDs := make([]string, len(filter.UserIDs))

		for i, v := range filter.UserIDs {
			uIDs[i] = v.String()
		}

		builder = builder.Where(sq.Eq{"id": uIDs})
	}

	query, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}

	err = r.db.SelectContext(ctx, &rawUsers, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("users not found")
		}

		return nil, err
	}

	return converter.ToDomainUserList(rawUsers), nil
}

func (r *UserRepository) Count(ctx context.Context, filter *domain.UserListFilter) (uint32, error) {
	builder := r.builder.Select("COUNT(id)").From(usersTable)

	if len(filter.UserIDs) != 0 {
		uIDs := make([]string, len(filter.UserIDs))

		for i, v := range filter.UserIDs {
			uIDs[i] = v.String()
		}

		builder = builder.Where(sq.Eq{"id": uIDs})
	}

	query, args, err := builder.ToSql()

	if err != nil {
		return 0, err
	}

	var total uint32

	err = r.db.GetContext(ctx, &total, query, args...)

	if err != nil {
		return 0, err
	}

	return total, nil
}
