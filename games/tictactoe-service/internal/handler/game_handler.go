package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pramanandasarkar02/game-server/games/tictactoe-service/internal/service"
)



type GameHandler struct{
	gameService *service.GameService
}


func NewGameHandler() *GameHandler{
	return &GameHandler{
		gameService: service.NewGameService(),
	}
}





func (h *GameHandler) CreateGame(c *gin.Context){
	game := h.gameService.CreateGame()
	c.JSON(http.StatusOK, game)
}


func (h *GameHandler)GetGame(c *gin.Context){
	id := c.Param("id")
	game, exists := h.gameService.GetGame(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
	}
	c.JSON(http.StatusOK, game)
}

func (h *GameHandler)GetAllGames(c *gin.Context){
	games := h.gameService.GetAllGames()
	c.JSON(http.StatusOK, games)
}




func (h *GameHandler)MakeMove(c *gin.Context){
	id := c.Param("id")
	var move struct{
		PositionX int `json:"positionX"`
		PositionY int `json:"positionY"`
		Player string `json:"player"`
	}
	if err := c.ShouldBindJSON(&move); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	game, err := h.gameService.MakeMove(id, move.PositionX, move.PositionY, move.Player)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, game)
}


