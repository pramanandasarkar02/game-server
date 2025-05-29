package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Player struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Level float32 `json:"level"`
}

type Game struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	RequiredPlayer int `json:"requiredPlayer"`
}

var (
	players []Player
	games   []Game
)

func playerConnection(c *gin.Context) {
	var newPlayer Player
	if err := c.ShouldBindJSON(&newPlayer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	players = append(players, newPlayer)

	log.Printf("%s is connected", newPlayer.Name)

	c.JSON(http.StatusOK, gin.H{
		"message": "Player connected successfully",
		"player":  newPlayer,
	})
}

func getPlayers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"players": players,
	})
}

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/connect", playerConnection)
	router.GET("/players", getPlayers)

	if err := router.Run(":4000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}