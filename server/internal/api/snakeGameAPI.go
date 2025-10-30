package api

import (
	"game-server/internal/handler"
	"game-server/internal/snake"
	"github.com/gin-gonic/gin"
)



func SnakeGameDataRoutes(router *gin.Engine) {
	// create snake service to communicate each other
	snakeService := snake.NewSnakeService()
	
	// create snake handler to handle snake game meta data
	snakeGameHandler := handler.NewSnakeHandler(snakeService)

	// get game specific meta data
	router.GET("/api/game/snake/meta-data",snakeGameHandler.MetaData)
	// get player match specific metadata
	router.GET("/api/game/snake/meta-data/:playerId", snakeGameHandler.GameMetaData)
	// main game logic end point 
	router.GET("/ws", snake.WsHandler)
}
