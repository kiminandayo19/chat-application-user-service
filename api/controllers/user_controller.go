package controllers

import "github.com/gin-gonic/gin"

type UserControllerInterface interface {
	GetAllUser(ctx *gin.Context)
}
