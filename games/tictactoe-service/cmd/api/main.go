package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pramanandasarkar02/game-server/games/tictactoe-service/internal/handler"
)


var (
	TICTACTOR_SERVICE_PORT = ":9000"
)


func main() {
	router := gin.Default()
	gameHandler := handler.NewGameHandler()

	router.POST("/", gameHandler.CreateGame)
	router.GET("/:id", gameHandler.GetGame)
	router.POST("/:id/move", gameHandler.MakeMove)


	log.Printf("Starting server on http://localhost%s", TICTACTOR_SERVICE_PORT)
	router.Run(TICTACTOR_SERVICE_PORT)
}