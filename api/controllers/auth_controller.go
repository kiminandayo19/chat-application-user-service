package controllers

import "github.com/gin-gonic/gin"

type AuthControllerInterface interface {
	RegisterController(ctx *gin.Context)
	LoginController(ctx *gin.Context)
	RefreshTokenController(ctx *gin.Context)
	ChangePasswordController(ctx *gin.Context)
	DeleteAccountController(ctx *gin.Context)
	LogoutController(ctx *gin.Context)
}
