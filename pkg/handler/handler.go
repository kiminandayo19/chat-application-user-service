package handler

import (
	"chat-application/users/internal/domain"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type serviceFn[T any] func(ctx context.Context, payload T) domain.APIBaseResponse

func JSONController[T any](service serviceFn[T]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload T

		if err := ctx.BindJSON(&payload); err != nil {
			response := domain.NewResponse(http.StatusBadRequest, false, "Invalid Body Request", nil)
			ctx.AbortWithStatusJSON(response.Code, response)
			return
		}

		validate := validator.New()
		if err := validate.Struct(&payload); err != nil {
			errValidation := err.(validator.ValidationErrors)
			response := domain.NewResponse(http.StatusBadRequest, false, errValidation[0].Error(), nil)
			ctx.AbortWithStatusJSON(response.Code, response)
			return
		}

		response := service(ctx.Request.Context(), payload)
		ctx.JSON(response.Code, response)
		return
	}
}

func QueryController[T any](service serviceFn[T]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload T

		if err := ctx.BindQuery(&payload); err != nil {
			response := domain.NewResponse(http.StatusBadRequest, false, "Invalid Query Params", nil)
			ctx.AbortWithStatusJSON(response.Code, response)
			return
		}

		validate := validator.New()
		if err := validate.Struct(&payload); err != nil {
			errValidation := err.(validator.ValidationErrors)
			response := domain.NewResponse(http.StatusBadRequest, false, errValidation[0].Error(), nil)
			ctx.AbortWithStatusJSON(response.Code, response)
			return
		}

		response := service(ctx.Request.Context(), payload)
		ctx.JSON(response.Code, response)
		return
	}
}
