package controllers

import (
	"QueryLib/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

func BookSeat(c *gin.Context) {
	fmt.Println("Hello World!")
	user := utils.NewSession()
	err := user.Login("user", "passwd")
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, gin.H{
		"message": err,
	})
}
