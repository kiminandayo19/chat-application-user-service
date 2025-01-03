package routes

import (
	"chat-application/users/api/controllers"
	cImpl "chat-application/users/api/controllers/impl"
	"chat-application/users/api/services"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	userController controllers.UserControllerInterface
}

func NewUserRoute(userService services.UserServiceInterface) *UserRoute {
	return &UserRoute{
		userController: cImpl.NewUserControllerImpl(userService),
	}
}

func (r *UserRoute) Setup(rg *gin.RouterGroup) {
	user := rg.Group("/user")
	{
		user.GET("/get-all", r.userController.GetAllUser)
	}
}
