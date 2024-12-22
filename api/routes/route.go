package routes

import (
	"chat-application/users/api/middlewares"
	"chat-application/users/config"
	"chat-application/users/internal/domain"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func welcomeHandler(ctx *gin.Context) {
	response := domain.NewResponse(
		http.StatusOK,
		true,
		"Chat Application - User Service",
		nil,
	)

	ctx.JSON(http.StatusOK, response)
}

func NewRoute() *gin.Engine {
	env, err := config.NewEnvDev()
	if err != nil {
		log.Panic("[router] - Failed to load env")
	}

	db := config.NewDbConnection()

	router := gin.Default()
	router.Use(middlewares.HeaderMiddleware(nil))

	basePath := "/v1/" + env.APP_ENTRY
	authRoute := NewAuthRoute(db, &env)

	v1 := router.Group(basePath)
	{
		v1.GET("/", welcomeHandler)
		authRoute.Setup(v1)
	}

	return router
}
