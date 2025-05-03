package main

import (
	"log"

	"github.com/gin-gonic/gin"
)


var (
	API_GATEWAY_PORT = ":8080"
)


func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})


	log.Printf("Starting server on http://localhost%s", API_GATEWAY_PORT)
	r.Run(API_GATEWAY_PORT)


}