package api

import (
	"game-server/internal/handler"
	"github.com/gin-gonic/gin"
)
// Match Make Routes for player Match Make 
func MatchMakeRoutes(router * gin.Engine){
	// create handler instances
	matchMakeHandler := handler.NewMatchMakeHandler()

	// Add to the queue
	router.POST("/api/match-make/:playerId/:gameId", matchMakeHandler.AddQueue)
	// Remove from the queue
	router.PATCH("/api/match-make/:playerId", matchMakeHandler.RemoveQueue)
	// Get Match if it already made
	router.GET("/api/match-make/:playerId", matchMakeHandler.GetMatch)
};