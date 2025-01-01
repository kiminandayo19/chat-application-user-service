package impl

import (
	"chat-application/users/api/controllers"
	"chat-application/users/api/services"
	"chat-application/users/internal/domain"
	"chat-application/users/internal/domain/dtos"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type authControllerImpl struct {
	authService services.AuthServiceInterface
}

func NewAuthControllerImpl(authService services.AuthServiceInterface) controllers.AuthControllerInterface {
	return &authControllerImpl{
		authService: authService,
	}
}

func (c *authControllerImpl) RegisterController(ctx *gin.Context) {
	var registerPayload dtos.RegisterRequestPayload

	if err := ctx.BindJSON(&registerPayload); err != nil {
		response := domain.NewResponse(http.StatusBadRequest, false, "Invalid Request Body", nil)
		ctx.AbortWithStatusJSON(response.Code, response)
		return
	}

	validate := validator.New()

	err := validate.Struct(registerPayload)
	if err != nil {
		errValidation := err.(validator.ValidationErrors)
		response := domain.NewResponse(http.StatusBadRequest, false, errValidation[0].Error(), nil)
		ctx.AbortWithStatusJSON(response.Code, response)
		return
	}

	response := c.authService.Register(ctx.Request.Context(), registerPayload)

	ctx.JSON(response.Code, response)
	return
}

func (c *authControllerImpl) LoginController(ctx *gin.Context) {
	var loginPayload dtos.LoginRequestPayload

	if err := ctx.BindJSON(&loginPayload); err != nil {
		response := domain.NewResponse(http.StatusBadRequest, false, "Invalid Request Body", nil)
		ctx.AbortWithStatusJSON(response.Code, response)
		return
	}

	validate := validator.New()
	err := validate.Struct(loginPayload)
	if err != nil {
		errValidation := err.(validator.ValidationErrors)
		response := domain.NewResponse(http.StatusBadRequest, false, errValidation[0].Error(), nil)
		ctx.AbortWithStatusJSON(response.Code, response)
		return
	}

	response := c.authService.Login(ctx.Request.Context(), loginPayload)

	ctx.JSON(response.Code, response)
	return
}

func (c *authControllerImpl) RefreshTokenController(ctx *gin.Context) {}

func (c *authControllerImpl) ChangePasswordController(ctx *gin.Context) {}

func (c *authControllerImpl) DeleteAccountController(ctx *gin.Context) {}
