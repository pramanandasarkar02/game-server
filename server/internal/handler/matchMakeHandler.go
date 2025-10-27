package handler

import (
	"fmt"
	"game-server/internal/service"
	"log"
	"github.com/gin-gonic/gin"
)

type MatchMakeHandler struct{
	matchMakeService *service.MatchMakeService
}

// Create new MatchMakeHandler
func NewMatchMakeHandler() *MatchMakeHandler{
	log.Println("Match-make service initiate")
	return &MatchMakeHandler{
		matchMakeService: service.NewMatchMakeService(),
	}
}
// Add player to the qeueue
func(mh *MatchMakeHandler)AddQueue(c *gin.Context){
	playerId := c.Param("playerId")
	gameId := c.Param("gameId")
	log.Printf("Player %v request for match-make for game %v", playerId, gameId)
	
	err := mh.matchMakeService.AddQueue(playerId, gameId)
	if err != nil{
		c.JSON(500, gin.H{
			"massage": "failed to add queue",
			"err": err,
		})
		return 
	}

	msg := fmt.Sprintf("%v add to the queue for the game %v", playerId, gameId)
	log.Println(msg)
	c.JSON(200, gin.H{
		"message": msg,
	})
}

// Remove Player from the queue
func(mh *MatchMakeHandler)RemoveQueue(c *gin.Context){
	playerId := c.Param("playerId")
	err := mh.matchMakeService.RemoveQueue(playerId)
	if err != nil{
		c.JSON(200, gin.H{
			"massage": "failed to remove queue",
			"err": err,
		})
		return 
	}

	msg := fmt.Sprintf("%v removed from the queue.", playerId)
	c.JSON(200, gin.H{
		"message": msg,

	})
}

// Get Match Stats from the queue
func (mh *MatchMakeHandler) GetMatch(c *gin.Context) {
	playerId := c.Param("playerId")

	match, err := mh.matchMakeService.GetMatch(playerId)
	if err != nil {
		c.JSON(500, gin.H{
			"isFound": false,
			"message": "failed to get match",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"isFound": true,
		"match":   match,
	})
}