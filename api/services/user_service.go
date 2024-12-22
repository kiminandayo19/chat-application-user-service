package services

import (
	"github.com/gin-gonic/gin"
)

type UserServiceInterface interface {
	Regist(ctx *gin.Context)
	Login(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	UpdatePassword(ctx *gin.Context)
	Delete(ctx *gin.Context)
}
