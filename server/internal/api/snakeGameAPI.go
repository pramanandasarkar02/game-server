package api

import (
	"game-server/internal/handler"
	"game-server/internal/snake"

	"github.com/gin-gonic/gin"
)



func SnakeGameDataRoutes(router *gin.Engine) {

	snakeService := snake.NewSnakeService()
	
	snakeGameHandler := handler.NewSnakeHandler(snakeService)

	router.GET("/api/game/snake/meta-data",snakeGameHandler.MetaData)
	router.GET("/api/game/snake/meta-data/:playerId", snakeGameHandler.GameMetaData)

	router.GET("ws", snake.WsHandler)

	
}
