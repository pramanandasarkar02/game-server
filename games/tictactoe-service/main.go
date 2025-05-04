package main


import (
	"log"

	"github.com/gin-gonic/gin"
)


var (
	TICTACTOR_SERVICE_PORT = ":9000"
)


func main() {
	router := gin.Default()


	log.Printf("Starting server on http://localhost%s", TICTACTOR_SERVICE_PORT)
	router.Run(TICTACTOR_SERVICE_PORT)
}