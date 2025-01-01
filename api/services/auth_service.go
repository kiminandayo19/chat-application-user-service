package services

import (
	"chat-application/users/internal/domain"
	dto "chat-application/users/internal/domain/dtos"
	"context"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, payload dto.RegisterRequestPayload) domain.APIBaseResponse
	Login(ctx context.Context, payload dto.LoginRequestPayload) domain.APIBaseResponse
	RefreshToken(ctx context.Context, payload dto.RefreshTokenRequestPayload)
	ChangePassword(ctx context.Context, payload dto.ForgotPasswordRequestPayload)
	Delete(ctx context.Context, payload dto.DeleteAccountRequestPayload)
}
