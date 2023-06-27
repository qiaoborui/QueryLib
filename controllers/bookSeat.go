package controllers

import (
	"QueryLib/services"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookSeatBody struct{
	Username string `json:"username"`
	Password  string `json:"password"`
	SeatID    string `json:"seatID"`
	Floor	 string `json:"floor"`
	Area	  string `json:"area"`
}
func BookSeat(c *gin.Context) {
	body,err := io.ReadAll(c.Request.Body)
	fmt.Println(body)
	var bookSeatBody BookSeatBody
	err = json.Unmarshal(body, &bookSeatBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error unmarshalling request body",
		})
		return
	}

	user := utils.NewSession()
	err = user.Login(bookSeatBody.Username, bookSeatBody.Password)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, gin.H{
		"message": err,
	})
}
