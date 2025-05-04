package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

var (
	USER_SERVICE_PORT = ":8082"
)


func main(){
	router := gin.Default()


	log.Printf("Starting server on http://localhost%s", USER_SERVICE_PORT)
	router.Run(USER_SERVICE_PORT)
}