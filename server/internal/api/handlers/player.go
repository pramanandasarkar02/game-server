package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pramanandasarkar02/game-server/internal/models"
	"github.com/pramanandasarkar02/game-server/internal/store"
	"github.com/pramanandasarkar02/game-server/pkg/logger"
)

func PlayerHandler(store *store.PlayerStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var player models.Player
		if err := c.ShouldBindJSON(&player); err != nil {
			logger.Error("Invalid player data: %v", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err := store.AddPlayer(player); err != nil {
			logger.Error("Failed to add player: %v", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		logger.Info("Player %s connected", player.Name)
		c.JSON(200, gin.H{"message": "Player connected successfully", "player": player})
	}
}

func GetPlayersHandler(store *store.PlayerStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		players := store.GetAll()
		c.JSON(200, gin.H{"players": players})
	}
}