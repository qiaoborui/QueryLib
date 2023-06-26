package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
)

func RequestInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Access the status we are sending
		status := c.Writer.Status()
		log.Println("Status:", status)

		// Log the client IP
		clientIP := c.ClientIP()
		log.Println("Client IP:", clientIP)

		// Log the user-agent
		userAgent := c.Request.UserAgent()
		log.Println("User-agent:", userAgent)
	}
}
