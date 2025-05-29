package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)


type Player struct{
	ID string
	Name string 
	Level float32
}

type Game struct{
	ID string
	Title string
	RequiredPlayer string
}





var(
	players []Player
	games []Game
)


func playerConnection(c *gin.Context){
	var newPlayer Player

	if err := c.ShouldBindBodyWithJSON(&newPlayer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	players = append(players, newPlayer)

	log.Println(newPlayer.Name + "is connected")

	c.JSON(http.StatusOK, gin.H{
		"message": "Player connected successfully",
		"player": newPlayer,
	})
}






func main(){
	router := gin.Default()
	router.GET("/ping", func(c * gin.Context){
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.POST("/connect", playerConnection)
	// router.

	router.Run(":4000")
}
