package converter

import (
	"github.com/google/uuid"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/auth/internal/domain"
	"github.com/javascriptizer1/grpc-cli-chat.mono/service/auth/internal/repository/dao"
)

func ToDaoUser(domainUser *domain.User) *dao.User {
	return &dao.User{
		ID:           domainUser.ID.String(),
		Name:         domainUser.Name,
		Email:        domainUser.Email,
		PasswordHash: domainUser.Password,
		Role:         dao.UserRole(domainUser.Role),
		CreatedAt:    domainUser.CreatedAt,
		UpdatedAt:    domainUser.UpdatedAt,
	}
}

func ToDomainUser(repoUser *dao.User) *domain.User {
	user, _ := domain.NewUserWithID(
		uuid.MustParse(repoUser.ID),
		repoUser.Name,
		repoUser.Email,
		repoUser.PasswordHash,
		domain.UserRole(repoUser.Role),
		repoUser.CreatedAt,
		repoUser.UpdatedAt,
	)

	return user
}

func ToDomainUserList(repoUsers []*dao.User) []*domain.User {
	domainUsers := make([]*domain.User, len(repoUsers))

	for i, repoUser := range repoUsers {
		domainUsers[i] = ToDomainUser(repoUser)
	}

	return domainUsers
}
