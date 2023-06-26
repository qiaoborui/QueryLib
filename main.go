package main

import (
	"QueryLib/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	fmt.Println("Starting server...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	host := os.Getenv("wzxy_HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	fmt.Printf("Server listening on http://%s:%s\n", host, port)
	r := gin.Default()
	routes.InitializeRoutes(r)
	r.Run(host + ":" + port)
}
