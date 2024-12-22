package routes

import (
	"chat-application/users/api/services"
	"chat-application/users/api/services/impl"
	"chat-application/users/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthRoute struct {
	userService services.UserServiceInterface
}

func NewAuthRoute(db *gorm.DB, env *config.EnvDevType) *AuthRoute {
	return &AuthRoute{
		userService: impl.NewUserServiceImpl(db, env),
	}
}

func (r *AuthRoute) Setup(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", r.userService.Regist)
		auth.POST("/login", r.userService.Login)
		auth.POST("/refresh-token", r.userService.RefreshToken)
		auth.POST("/update-password", r.userService.UpdatePassword)
		auth.DELETE("/delete-user", r.userService.Delete)
	}
}
