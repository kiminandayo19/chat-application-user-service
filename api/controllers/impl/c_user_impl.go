package impl

import (
	"chat-application/users/api/controllers"
	"chat-application/users/api/services"
	"chat-application/users/internal/domain"
	"chat-application/users/internal/domain/dtos"
	"chat-application/users/pkg/handler"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type userControllerImpl struct {
	userService services.UserServiceInterface
}

func NewUserControllerImpl(userService services.UserServiceInterface) controllers.UserControllerInterface {
	return &userControllerImpl{
		userService: userService,
	}
}

func (c *userControllerImpl) GetAllUser(ctx *gin.Context) {
	query := ctx.Request.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	payload := dtos.PaginatedRequestPayload{
		Page:  page,
		Limit: limit,
	}

	if payload.Page <= 0 || payload.Limit < 10 {
		response := domain.NewResponse(http.StatusBadRequest, false, "Invalid Query Params", nil)
		ctx.AbortWithStatusJSON(response.Code, response)
		return
	}
	handler.QueryController(c.userService.Get, payload)(ctx)
	return
}
