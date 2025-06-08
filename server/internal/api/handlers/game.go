package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pramanandasarkar02/game-server/internal/store"
)

func GetGamesHandler(store store.GameStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		games := store.GetAll()
		gamesCopy := make([]interface{}, len(games))
		for i, g := range games {
			gamesCopy[i] = struct {
				ID             string `json:"id"`
				Title          string `json:"title"`
				RequiredPlayer int    `json:"requiredPlayer"`
			}{
				ID:             g.ID(),
				Title:          g.Title(),
				RequiredPlayer: g.RequiredPlayers(),
			}
		}
		c.JSON(200, gin.H{"games": gamesCopy})
	}
}