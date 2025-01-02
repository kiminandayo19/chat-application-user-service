package routes

import (
	"chat-application/users/api/controllers"
	cImpl "chat-application/users/api/controllers/impl"
	"chat-application/users/api/services"

	"github.com/gin-gonic/gin"
)

type AuthRoute struct {
	userController controllers.AuthControllerInterface
}

func NewAuthRoute(authService services.AuthServiceInterface) *AuthRoute {
	return &AuthRoute{
		userController: cImpl.NewAuthControllerImpl(authService),
	}
}

func (r *AuthRoute) Setup(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", r.userController.RegisterController)
		auth.POST("/login", r.userController.LoginController)
		auth.POST("/refresh-token", r.userController.RefreshTokenController)
		// auth.POST("/update-password", r.userController.ChangePasswordController)
		// auth.DELETE("/delete-account", r.userController.DeleteAccountController)
	}
}
