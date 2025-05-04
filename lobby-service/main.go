package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

var (
	LOBBY_SERVICE_PORT = ":8081"
)


func main(){
	router := gin.Default()



	log.Printf("Starting server on http://localhost%s", LOBBY_SERVICE_PORT)
	router.Run(LOBBY_SERVICE_PORT)

}