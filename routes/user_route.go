package routes

import (
	"gin-mongo-api/controllers" //add this

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.POST("/user", controllers.CreateUser()) //add this
}