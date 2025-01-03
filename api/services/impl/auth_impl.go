package impl

import (
	"chat-application/users/api/services"
	"chat-application/users/config"
	"chat-application/users/internal/domain"
	"chat-application/users/internal/domain/dtos"
	dto "chat-application/users/internal/domain/dtos"
	"chat-application/users/internal/domain/models"
	repo "chat-application/users/internal/domain/repositories"
	"chat-application/users/pkg/jwt"
	"chat-application/users/pkg/utils"
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	jwts "github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type authServiceImpl struct {
	repo repo.AuthRepository
	rdb  *redis.Client
}

func NewAuthService(repo repo.AuthRepository, rdb *redis.Client) services.AuthServiceInterface {
	return &authServiceImpl{
		repo: repo,
		rdb:  rdb,
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

	redisErr := s.rdb.Set(ctx, "access-"+utilsHelper.UintToString(find.ID), accessToken, 15*time.Minute).Err()
	redisErr = s.rdb.Set(ctx, "refresh-"+utilsHelper.UintToString(find.ID), refreshToken, 24*7*time.Hour).Err()
	if redisErr != nil {
		return domain.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
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

func (s *authServiceImpl) RefreshToken(ctx context.Context, payload dto.RefreshTokenRequestPayload) domain.APIBaseResponse {
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

	decodeData, err := refreshHelper.ValidateToken(payload.RefreshToken)
	err = s.rdb.Get(ctx, "refresh-"+decodeData.UserId).Err()
	if err != nil {
		switch {
		case errors.Is(err, jwts.ErrTokenExpired):
			return domain.NewResponse(http.StatusUnauthorized, false, "Token Expired", nil)
		case errors.Is(err, jwts.ErrTokenInvalidId):
			return domain.NewResponse(http.StatusBadRequest, false, "Invalid Token", nil)
		default:
			return domain.NewResponse(http.StatusInternalServerError, false, "Failed to validate token", nil)
		}
	}

	user := &models.Mst_users{
		Username: decodeData.Username,
	}

	find, err := s.repo.FindByUsername(ctx, user)

	if err != nil {
		code := http.StatusInternalServerError

		if err.Error() == repo.ErrRecordNotFound {
			code = http.StatusForbidden
		}
		return domain.NewResponse(code, false, err.Error(), nil)
	}

	accessToken, err := accessHelper.GenerateToken(utilsHelper.UintToString(find.ID), find.Username)
	if err != nil {
		return domain.NewResponse(http.StatusInternalServerError, false, "Internal Server Error. Failed to generate jwt", nil)
	}

	err = s.rdb.Set(ctx, "access-"+utilsHelper.UintToString(find.ID), accessToken, 0).Err()
	if err != nil {
		return domain.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
	}

	response := dtos.RefreshTokenResponsePayload{
		AccessToken: accessToken,
	}
	return domain.NewResponse(http.StatusOK, true, "Success", response)
}

func (s *authServiceImpl) ChangePassword(ctx context.Context, payload dto.ChangePasswordRequestPayload) domain.APIBaseResponse {
	if ctx.Err() != nil {
		return domain.NewResponse(http.StatusBadRequest, false, "Bad Request", nil)
	}

	utilsHelper := utils.NewUtils()

	user := &models.Mst_users{
		ID: uint(payload.UserId.(float64)),
	}

	find, err := s.repo.FindByID(ctx, user)
	if err != nil {
		code := http.StatusInternalServerError
		if err.Error() == repo.ErrRecordNotFound {
			code = http.StatusBadRequest
		}
		return domain.NewResponse(code, false, err.Error(), nil)
	}

	if err := utilsHelper.ComparePassword(find.Password, payload.Password); err == nil {
		return domain.NewResponse(http.StatusBadRequest, false, "New Password is the same as old password, choose different one.", nil)
	}

	hashPassword, err := utilsHelper.HashPassword(payload.Password)
	if err != nil {
		return domain.NewResponse(http.StatusInternalServerError, false, "Failed to hash password", nil)
	}

	if err := s.repo.UpdateById(ctx, find.ID, "password", hashPassword); err != nil {
		code := http.StatusInternalServerError
		if err.Error() == repo.ErrRecordNotFound {
			code = http.StatusBadRequest
		}
		return domain.NewResponse(code, false, err.Error(), nil)
	}

	return domain.NewResponse(http.StatusOK, true, "Success Update Password", nil)
}

func (s *authServiceImpl) Delete(ctx context.Context, payload dto.DeleteAccountRequestPayload) domain.APIBaseResponse {
	if ctx.Err() != nil {
		return domain.NewResponse(http.StatusBadRequest, false, "Bad Request", nil)
	}
	userId, _ := strconv.Atoi(payload.UserId.(string))

	if err := s.repo.UpdateById(ctx, uint(userId), "is_deleted", true); err != nil {
		code := http.StatusInternalServerError
		message := err.Error()
		switch {
		case err.Error() == repo.ErrRecordNotFound:
			code = http.StatusBadRequest
		case err.Error() == repo.ErrNotUpdated:
			code = http.StatusBadRequest
			message = "Account already deleted."
		}
		return domain.NewResponse(code, false, message, nil)
	}

	return domain.NewResponse(http.StatusOK, true, "Success Delete Account", nil)
}

func (s *authServiceImpl) Logout(ctx context.Context, payload dto.LogoutRequestPayload) domain.APIBaseResponse {
	if ctx.Err() != nil {
		return domain.NewResponse(http.StatusBadRequest, false, "Bad Request", nil)
	}
	userId, _ := strconv.Atoi(payload.UserId.(string))

	user := &models.Mst_users{
		ID: uint(userId),
	}

	find, err := s.repo.FindByID(ctx, user)
	if err != nil {
		code := http.StatusInternalServerError
		if err.Error() == repo.ErrRecordNotFound {
			code = http.StatusBadRequest
		}
		return domain.NewResponse(code, false, err.Error(), nil)
	}

	rdbErr := s.rdb.GetDel(ctx, "access-"+utils.NewUtils().UintToString(find.ID)).Err()
	rdbErr = s.rdb.GetDel(ctx, "refresh-"+utils.NewUtils().UintToString(find.ID)).Err()
	if rdbErr != nil {
		return domain.NewResponse(http.StatusInternalServerError, false, "Internal Server Error", nil)
	}

	return domain.NewResponse(http.StatusOK, true, "Success", nil)
}
