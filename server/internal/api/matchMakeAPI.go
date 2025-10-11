package api

import (
	"game-server/internal/handler"

	"github.com/gin-gonic/gin"
)

func MatchMakeRoutes(router * gin.Engine){
	matchMakeHandler := handler.NewMatchMakeHandler()
	
	router.POST("/api/match-make/:playerId/:gameId", matchMakeHandler.AddQueue)
	router.PATCH("/api/match-make/:playerId", matchMakeHandler.RemoveQueue)
	router.GET("/api/match-make/:playerId", matchMakeHandler.GetMatch)
}