package routes

import (
	"chat-application/users/api/controllers"
	cImpl "chat-application/users/api/controllers/impl"
	"chat-application/users/api/services"

	"github.com/gin-gonic/gin"
)

type AuthRoute struct {
	authController controllers.AuthControllerInterface
}

func NewAuthRoute(authService services.AuthServiceInterface) *AuthRoute {
	return &AuthRoute{
		authController: cImpl.NewAuthControllerImpl(authService),
	}
}

func (r *AuthRoute) Setup(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", r.authController.RegisterController)
		auth.POST("/login", r.authController.LoginController)
		auth.POST("/refresh-token", r.authController.RefreshTokenController)
		auth.POST("/change-password", r.authController.ChangePasswordController)
		auth.DELETE("/delete-account", r.authController.DeleteAccountController)
	}
}
