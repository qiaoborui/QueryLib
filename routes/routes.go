package routes

import (
	"QueryLib/controllers"
	"QueryLib/middlewares"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {
	router.Use(middlewares.RequestInfo())
	router.POST("/queryUsers", controllers.QueryUsers)
	router.GET("/queryUser", controllers.QueryUser)
	router.GET("/bookSeat", controllers.BookSeat)
}
