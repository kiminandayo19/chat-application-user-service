package services

import (
	"chat-application/users/internal/domain"
	dto "chat-application/users/internal/domain/dtos"
	"context"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, payload dto.RegisterRequestPayload) domain.APIBaseResponse
	Login(ctx context.Context, payload dto.LoginRequestPayload) domain.APIBaseResponse
	RefreshToken(ctx context.Context, payload dto.RefreshTokenRequestPayload) domain.APIBaseResponse
	ChangePassword(ctx context.Context, payload dto.ChangePasswordRequestPayload) domain.APIBaseResponse
	Delete(ctx context.Context, payload dto.DeleteAccountRequestPayload) domain.APIBaseResponse
	Logout(ctx context.Context, payload dto.LogoutRequestPayload) domain.APIBaseResponse
}
