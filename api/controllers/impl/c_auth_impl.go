package impl

import (
	"chat-application/users/api/controllers"
	"chat-application/users/api/services"
	"chat-application/users/internal/domain"
	"chat-application/users/internal/domain/dtos"
	"chat-application/users/pkg/handler"
	"net/http"

	"github.com/gin-gonic/gin"
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
	handler.JSONController(c.authService.Register)(ctx)
	return
}

func (c *authControllerImpl) LoginController(ctx *gin.Context) {
	// equivalent to handler := ...; handler(ctx)
	handler.JSONController(c.authService.Login)(ctx)
	return
}

func (c *authControllerImpl) RefreshTokenController(ctx *gin.Context) {
	handler.JSONController(c.authService.RefreshToken)(ctx)
	return
}

func (c *authControllerImpl) ChangePasswordController(ctx *gin.Context) {
	handler.JSONController(c.authService.ChangePassword)(ctx)
	return
}

func (c *authControllerImpl) DeleteAccountController(ctx *gin.Context) {
	payload := dtos.DeleteAccountRequestPayload{
		UserId: ctx.Request.URL.Query().Get("userId"),
	}

	if payload.UserId == nil {
		response := domain.NewResponse(http.StatusBadRequest, false, "Invalid Query Params", nil)
		ctx.AbortWithStatusJSON(response.Code, response)
	}

	handler.QueryController(c.authService.Delete, payload)(ctx)
	return
}
