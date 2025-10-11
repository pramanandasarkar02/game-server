package handler

import (
	"fmt"
	"game-server/internal/service"

	"github.com/gin-gonic/gin"
)
type MatchMakeHandler struct{
	matchMakeService *service.MatchMakeService
}



func NewMatchMakeHandler() *MatchMakeHandler{
	return &MatchMakeHandler{
		matchMakeService: service.NewMatchMakeService(),
	}
}


func(mh *MatchMakeHandler)AddQueue(c *gin.Context){
	playerId := c.Param("playerId")
	gameId := c.Param("gameId")

	err := mh.matchMakeService.AddQueue(playerId, gameId)
	if err != nil{
		c.JSON(200, gin.H{
			"massage": "failed to add queue",
			"err": err,
		})
		return 
	}

	msg := fmt.Sprintf("%v add to the queue for the game %v", playerId, gameId)
	c.JSON(200, gin.H{
		"message": msg,

	})
}


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