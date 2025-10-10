package api

import (
	"game-server/internal/handler"

	"github.com/gin-gonic/gin"
)



func SnakeGameDataRoutes(router *gin.Engine) {
	
	snakeGameHandler := handler.NewSnakeHandler()

	router.GET("/api/game/snake/meta-data",snakeGameHandler.MetaData)
	router.GET("/api/game/snake/meta-data/:playerId", snakeGameHandler.GameMetaData)


	
}
