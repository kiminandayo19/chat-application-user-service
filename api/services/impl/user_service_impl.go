package impl

import (
	"chat-application/users/api/services"
	"chat-application/users/config"
	"chat-application/users/internal/domain"
	dto "chat-application/users/internal/domain/dtos"
	"chat-application/users/internal/domain/models"
	"chat-application/users/pkg/jwt"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	DB  *gorm.DB
	env *config.EnvDevType
}

func NewUserServiceImpl(db *gorm.DB, env *config.EnvDevType) services.UserServiceInterface {
	return &UserServiceImpl{
		DB:  db,
		env: env,
	}
}

func (s *UserServiceImpl) Regist(ctx *gin.Context) {
	var user dto.RegisterRequestPayload
	var data models.Mst_users

	if err := ctx.BindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, domain.NewResponse(http.StatusBadRequest, false, "Invalid request body", nil))
		return
	}

	validate := validator.New()

	if err := validate.Struct(&user); err != nil {
		errValidation := err.(validator.ValidationErrors)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, domain.NewResponse(http.StatusBadRequest, false, errValidation[0].Error(), nil))
		return
	}

	byte, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, domain.NewResponse(http.StatusInternalServerError, false, "Failed to hash password", nil))
		return
	}

	data.Username = user.Username
	data.Phonenumber = user.Phonenumber
	data.Email = user.Email
	data.Password = string(byte)

	if err := s.DB.Create(&data).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, domain.NewResponse(http.StatusInternalServerError, false, err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusCreated, domain.NewResponse(http.StatusCreated, true, "Register success", nil))
	return
}

func (s *UserServiceImpl) Login(ctx *gin.Context) {
	var user dto.LoginRequestPayload
	var storedUser models.Mst_users
	env := s.env

	if err := ctx.BindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, domain.NewResponse(http.StatusBadRequest, false, "Invalid request body", nil))
		return
	}

	validate := validator.New()

	if err := validate.Struct(&user); err != nil {
		errValidation := err.(validator.ValidationErrors)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, domain.NewResponse(http.StatusBadRequest, false, errValidation[0].Error(), nil))
		return
	}

	find := s.DB.Select("id", "password", "username", "email", "phonenumber").Where(&models.Mst_users{Username: user.Username}).First(&storedUser)

	if find.Error != nil {
		if errors.Is(find.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, domain.NewResponse(http.StatusBadRequest, false, "User not found", nil))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, domain.NewResponse(http.StatusInternalServerError, false, find.Error.Error(), nil))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, domain.NewResponse(http.StatusBadRequest, false, "Wrong password", nil))
		return
	}

	accessJwtHelper := jwt.NewJWTHelper(env.JWT_SECRET, 15*time.Minute)
	refreshJwtHelper := jwt.NewJWTHelper(env.JWT_SECRET, 24*time.Hour)

	accessToken, err := accessJwtHelper.GenerateToken(string(strconv.FormatUint(uint64(storedUser.ID), 10)), storedUser.Username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, domain.NewResponse(http.StatusInternalServerError, false, err.Error(), nil))
		return
	}

	refreshToken, err := refreshJwtHelper.GenerateToken(string(strconv.FormatUint(uint64(storedUser.ID), 10)), storedUser.Username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, domain.NewResponse(http.StatusInternalServerError, false, err.Error(), nil))
		return
	}

	response := dto.LoginResponsePayload{
		UserId:       string(strconv.FormatUint(uint64(storedUser.ID), 10)),
		Username:     storedUser.Username,
		Email:        storedUser.Email,
		Phonenumber:  storedUser.Phonenumber,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	ctx.JSON(http.StatusOK, domain.NewResponse(http.StatusOK, true, "Login success", response))
	return
}

func (s *UserServiceImpl) RefreshToken(ctx *gin.Context) {
	var refreshToken dto.RefreshTokenRequestPayload
	env := s.env

	if err := ctx.BindJSON(&refreshToken); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, domain.NewResponse(http.StatusBadRequest, false, "Invalid body", nil))
		return
	}

	accessJwtHelper := jwt.NewJWTHelper(env.JWT_SECRET, 15*time.Minute)

	newToken, err := accessJwtHelper.RefreshToken(refreshToken.RefreshToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, domain.NewResponse(http.StatusInternalServerError, false, err.Error(), nil))
		return
	}

	response := dto.RefreshTokenResponsePayload{
		AccessToken: newToken,
	}

	ctx.JSON(http.StatusOK, domain.NewResponse(http.StatusOK, true, "Success", response))
	return
}

func (s *UserServiceImpl) UpdatePassword(ctx *gin.Context) {}

func (s *UserServiceImpl) Delete(ctx *gin.Context) {}
