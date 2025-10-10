package handler

import (
	"game-server/internal/snake"

	"github.com/gin-gonic/gin"
)


type SnakeHandler struct{
	snakeService *snake.SnakeService
}


func NewSnakeHandler(ss * snake.SnakeService) *SnakeHandler{
	return &SnakeHandler{
		snakeService: ss,
	}
}


func (sh *SnakeHandler) MetaData(c *gin.Context) {
	gameMetaData := sh.snakeService.SnakeGameMetaData()
	c.JSON(200, gameMetaData)
}

func (sh *SnakeHandler) GameMetaData(c *gin.Context){
	playerId := c.Param("playerId")

	c.JSON(200, gin.H{
		"message": "player specific game meta data",
		"playerId": playerId,
	})
}