package middlewares

import (
	"chat-application/users/internal/domain"
	"chat-application/users/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func abortHandler(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusUnauthorized, domain.NewResponse(http.StatusUnauthorized, false, message, nil))
	ctx.Abort()
	return
}

func AuthMiddleware(jwtHelper *jwt.JWTHelper) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			abortHandler(ctx, "Unauthorized")
			return
		}

		token, err := jwt.ExtractBearerToken(authHeader)
		if err != nil {
			abortHandler(ctx, err.Error())
			return
		}

		claims, err := jwtHelper.ValidateToken(token)
		if err != nil {
			abortHandler(ctx, err.Error())
			return
		}

		ctx.Set("userId", claims.UserId)
		ctx.Set("username", claims.Username)
		ctx.Next()
	}
}
