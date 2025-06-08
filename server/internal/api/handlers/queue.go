package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pramanandasarkar02/game-server/internal/store"
	"github.com/pramanandasarkar02/game-server/pkg/logger"
)

func GetQueueHandler(store store.QueueStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		queues := store.GetAllQueues()
		c.JSON(200, gin.H{"queue": queues})
	}
}

func JoinQueueHandler(playerStore store.PlayerStore, gameStore store.GameStore, queueStore store.QueueStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			PlayerID string `json:"playerID"`
			GameID   string `json:"gameID"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error("Invalid queue join request: %v", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		_, exists := playerStore.GetPlayer(req.PlayerID)
		if !exists {
			c.JSON(404, gin.H{"message": fmt.Sprintf("player with ID '%s' not found", req.PlayerID)})
			return
		}

		_, exists = gameStore.GetGame(req.GameID)
		if !exists {
			c.JSON(404, gin.H{"message": fmt.Sprintf("game with ID '%s' not found", req.GameID)})
			return
		}

		queueStore.AddPlayer(req.GameID, req.PlayerID)
		logger.Info("Player %s joined queue for game %s", req.PlayerID, req.GameID)
		c.JSON(200, gin.H{"message": fmt.Sprintf("player(%s) join the queue for game %s", req.PlayerID, req.GameID)})
	}
}