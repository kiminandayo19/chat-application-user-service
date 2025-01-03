package impl

import (
	"chat-application/users/api/services"
	"chat-application/users/internal/domain"
	"chat-application/users/internal/domain/dtos"
	repo "chat-application/users/internal/domain/repositories"
	"context"
	"net/http"
)

type userServiceImpl struct {
	repo repo.UserRepository
}

func NewUserService(repo repo.UserRepository) services.UserServiceInterface {
	return &userServiceImpl{
		repo: repo,
	}
}

func (s *userServiceImpl) Get(ctx context.Context, payload dtos.PaginatedRequestPayload) domain.APIBaseResponse {
	if ctx.Err() != nil {
		return domain.NewResponse(http.StatusBadRequest, false, "Bad Request", nil)
	}

	users, total, err := s.repo.GetAll(ctx, payload.Page, payload.Limit)
	if err != nil {
		return domain.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
	}

	return domain.NewResponse(http.StatusOK, true, "Success", users).NewResponseWithPagination(payload.Page, payload.Limit, total)
}
