package impl

import (
	"chat-application/users/api/services"
	"chat-application/users/config"
	"chat-application/users/internal/domain"
	dto "chat-application/users/internal/domain/dtos"
	"chat-application/users/internal/domain/models"
	repo "chat-application/users/internal/domain/repositories"
	"chat-application/users/pkg/jwt"
	"chat-application/users/pkg/utils"
	"context"
	"log"
	"net/http"
	"time"
)

type authServiceImpl struct {
	repo repo.AuthRepository
}

func NewAuthService(repo repo.AuthRepository) services.AuthServiceInterface {
	return &authServiceImpl{
		repo: repo,
	}
}

func (s *authServiceImpl) Register(ctx context.Context, payload dto.RegisterRequestPayload) domain.APIBaseResponse {
	if ctx.Err() != nil {
		return domain.NewResponse(http.StatusBadRequest, false, "Bad Request", nil)
	}

	hashedPassword, err := utils.NewUtils().HashPassword(payload.Password)
	if err != nil {
		return domain.NewResponse(http.StatusInternalServerError, false, err.Error(), nil)
	}

	user := &models.Mst_users{
		Username:    payload.Username,
		Email:       payload.Email,
		Phonenumber: payload.Phonenumber,
		Password:    hashedPassword,
	}

	if err := s.repo.Insert(ctx, user); err != nil {
		code := http.StatusInternalServerError

		if err.Error() == repo.ErrDuplicatedKey {
			code = http.StatusBadRequest
		}
		return domain.NewResponse(code, false, err.Error(), nil)
	}

	return domain.NewResponse(http.StatusCreated, true, "Register success", nil)
}

func (s *authServiceImpl) Login(ctx context.Context, payload dto.LoginRequestPayload) domain.APIBaseResponse {
	if ctx.Err() != nil {
		return domain.NewResponse(http.StatusBadRequest, false, "Bad Request", nil)
	}

	env, err := config.NewEnvDev()
	if err != nil {
		log.Fatal("[server] - failed to load env")
	}

	utilsHelper := utils.NewUtils()
	accessHelper := jwt.NewJWTHelper(env.JWT_SECRET, 15*time.Minute)
	refreshHelper := jwt.NewJWTHelper(env.JWT_SECRET, 24*7*time.Hour)

	user := &models.Mst_users{
		Username: payload.Username,
		Password: payload.Password,
	}

	find, err := s.repo.FindByUsername(ctx, user)

	if err != nil {
		code := http.StatusInternalServerError

		if err.Error() == repo.ErrRecordNotFound {
			code = http.StatusBadRequest
		}
		return domain.NewResponse(code, false, err.Error(), nil)
	}

	if err := utilsHelper.ComparePassword(find.Password, payload.Password); err != nil {
		log.Println(err.Error())
		return domain.NewResponse(http.StatusBadRequest, false, "Invalid Password", nil)
	}

	accessToken, err := accessHelper.GenerateToken(utilsHelper.UintToString(find.ID), find.Username)
	refreshToken, err := refreshHelper.GenerateToken(utilsHelper.UintToString(find.ID), find.Username)

	if err != nil {
		return domain.NewResponse(http.StatusInternalServerError, false, "Internal Server Error. Failed to generate jwt", nil)
	}

	response := dto.LoginResponsePayload{
		UserId:       utilsHelper.UintToString(find.ID),
		Username:     find.Username,
		Email:        find.Email,
		Phonenumber:  find.Phonenumber,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return domain.NewResponse(http.StatusOK, true, "Login Success", response)
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, payload dto.RefreshTokenRequestPayload) {
}

func (s *authServiceImpl) ChangePassword(ctx context.Context, payload dto.ForgotPasswordRequestPayload) {
}

func (s *authServiceImpl) Delete(ctx context.Context, payload dto.DeleteAccountRequestPayload) {
}
