package handler

import (
	"game-server/internal/snake"

	"github.com/gin-gonic/gin"
)


type SnakeHandler struct{
	snakeService *snake.SnakeBoard
}


func NewSnakeHandler() *SnakeHandler{
	return &SnakeHandler{

	}
}


func (sh *SnakeHandler) MetaData(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "snake game meta data",

	})
}

func (sh *SnakeHandler) GameMetaData(c *gin.Context){
	playerId := c.Param("playerId")

	c.JSON(200, gin.H{
		"message": "player specific game meta data",
		"playerId": playerId,
	})
}