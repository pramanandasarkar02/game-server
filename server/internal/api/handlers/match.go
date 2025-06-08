package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pramanandasarkar02/game-server/internal/store"
)

func GetMatchForUserHandler(store *store.MatchStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userId")
		match, exists := store.GetMatchForPlayer(userID)
		if !exists {
			c.JSON(404, gin.H{"message": "player is not in running match"})
			return
		}
		c.JSON(200, gin.H{"matchID": match.ID})
	}
}

func GetMatchHandler(store *store.MatchStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		matchID := c.Param("matchId")
		playerID := c.Param("userId")
		match, exists := store.GetMatch(matchID)
		if !exists {
			c.JSON(404, gin.H{"message": fmt.Sprintf("match %s not found", matchID)})
			return
		}
		userInMatch := false
		for _, pID := range match.Players {
			if pID == playerID {
				userInMatch = true
				break
			}
		}
		if !userInMatch {
			c.JSON(403, gin.H{"message": fmt.Sprintf("player %s is not in match %s", playerID, matchID)})
			return
		}
		c.JSON(200, gin.H{"matchID": matchID, "players": match.Players})
	}
}