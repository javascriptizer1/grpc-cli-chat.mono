package service

import (
	"context"
	"errors"
	"strconv"

	"time"

	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.backend/pkg/helper/jwt"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/auth/internal/domain"
	"github.com/javascriptizer1/grpc-cli-chat.backend/service/auth/internal/service/dto"
)

type AuthConfig struct {
	AccessTokenSecretKey  string
	AccessTokenDuration   time.Duration
	RefreshTokenSecretKey string
	RefreshTokenDuration  time.Duration
}

type AuthService struct {
	userRepo UserRepository
	config   AuthConfig
}

func NewAuthService(userRepo UserRepository, cfg AuthConfig) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, input dto.RegisterInputDto) (id uuid.UUID, err error) {
	v, _ := s.userRepo.OneByEmail(ctx, input.Email)

	if v != nil {
		return id, errors.New("user with this email already exists")
	}

	if input.Password != input.PasswordConfirm {
		return id, errors.New("passwords don`t match")
	}

	u, err := domain.NewUser(input.Name, input.Email, input.Password, input.Role)

	if err != nil {
		return id, err
	}

	err = u.HashPassword()

	if err != nil {
		return id, err
	}

	if err = s.userRepo.Create(ctx, u); err != nil {
		return id, err
	}

	return u.ID, nil
}

func (s *AuthService) Login(ctx context.Context, login string, password string) (string, error) {
	v, err := s.userRepo.OneByEmail(ctx, login)

	if v == nil {
		return "", errors.New("invalid login or password")
	}

	if err != nil {
		return "", err
	}

	u, err := domain.NewUserWithID(v.ID, v.Name, v.Email, v.Password, v.Role, v.CreatedAt, v.UpdatedAt)

	if err != nil {
		return "", err
	}

	if !u.CheckPassword(password) {
		return "", errors.New("invalid login or password")
	}

	t, err := jwt.GenerateToken(jwt.UserClaims{
		ID:   u.ID.String(),
		Role: strconv.Itoa(int(u.Role)),
	},
		s.config.RefreshTokenSecretKey,
		s.config.RefreshTokenDuration,
	)

	return t, err
}

func (s *AuthService) AccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := jwt.VerifyToken(refreshToken, s.config.RefreshTokenSecretKey)

	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	stringID, err := claims.GetSubject()

	if err != nil {
		return "", err
	}

	id, err := uuid.Parse(stringID)

	if err != nil {
		return "", err
	}

	u, err := s.userRepo.OneByID(ctx, uuid.UUID(id))

	if err != nil {
		return "", err
	}

	accessToken, err := jwt.GenerateToken(jwt.UserClaims{
		ID:   u.ID.String(),
		Role: strconv.Itoa(int(u.Role)),
	},
		s.config.AccessTokenSecretKey,
		s.config.AccessTokenDuration,
	)

	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, oldRefreshToken string) (string, error) {
	claims, err := jwt.VerifyToken(oldRefreshToken, s.config.RefreshTokenSecretKey)

	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	stringID, err := claims.GetSubject()

	if err != nil {
		return "", err
	}

	id, err := uuid.Parse(stringID)

	if err != nil {
		return "", err
	}

	u, err := s.userRepo.OneByID(ctx, uuid.UUID(id))

	if err != nil {
		return "", err
	}

	refreshToken, err := jwt.GenerateToken(jwt.UserClaims{
		ID:   u.ID.String(),
		Role: strconv.Itoa(int(u.Role)),
	},
		s.config.RefreshTokenSecretKey,
		s.config.RefreshTokenDuration,
	)

	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (s *AuthService) Check(_ context.Context, endpoint string, role domain.UserRole) bool {
	return endpoint != "" && role != 0
}
