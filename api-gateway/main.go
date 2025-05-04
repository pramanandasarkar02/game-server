package main

import (
	"log"

	"github.com/gin-gonic/gin"
)


var (
	API_GATEWAY_PORT = ":8080"
)


func main() {
	router := gin.Default()

	router.GET("/", func(context * gin.Context){
		context.JSON(200, gin.H{
			"message": "game server api gateway",
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})


	log.Printf("Starting server on http://localhost%s", API_GATEWAY_PORT)
	router.Run(API_GATEWAY_PORT)


}