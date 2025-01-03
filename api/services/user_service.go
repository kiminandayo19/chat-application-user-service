package services

import (
	"chat-application/users/internal/domain"
	"chat-application/users/internal/domain/dtos"
	"context"
)

type UserServiceInterface interface {
	Get(ctx context.Context, payload dtos.PaginatedRequestPayload) domain.APIBaseResponse
}
