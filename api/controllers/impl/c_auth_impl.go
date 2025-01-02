package impl

import (
	"chat-application/users/api/controllers"
	"chat-application/users/api/services"
	"chat-application/users/pkg/handler"

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

// func (c *authControllerImpl) ChangePasswordController(ctx *gin.Context) {
// 	handler.JSONController(c.authService.ChangePassword)(ctx)
// 	return
// }
//
// func (c *authControllerImpl) DeleteAccountController(ctx *gin.Context) {
// 	handler.QueryController(c.authService.Delete)(ctx)
// 	return
// }
