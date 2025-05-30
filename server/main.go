package main

import (
	"fmt"
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
	queue map[interface{}]interface{}
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

func getGames(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"games": games,
	})
}

func getQueue(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"queue": queue,
	})
}

func enterQueue(c *gin.Context){
	var player Player
	var game Game
	type Request struct{
		playerId string
		gameId string
	}
	var req Request
	// find the player in server and then add him to the lobby 
	if err := c.ShouldBindBodyWithJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	foundPlayer := false

	for _, p := range players{
		if p.ID == req.playerId {
			player = p
			foundPlayer = true
		}
	}

	if !foundPlayer {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("player with ID '%s' not found", req.playerId)})
		return 
	}


	foundGame := false

	for _, g := range games{
		if g.ID == req.gameId {
			game = g
			foundGame = true
		}
	}

	if !foundGame {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("game with ID '%s' not found", req.gameId)})
		return 
	}

	

	queue[player.ID] = game.ID

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("player(%v) join the queue for game %v(%v)", req.playerId, game.Title, req.gameId)})


}


func createGames(){
	games = append(games, Game{
		ID: "a2",
		Title: "Tic Tac Toe",
		RequiredPlayer: 2,
	})
	games = append(games, Game{
		ID: "a4",
		Title: "Ludu",
		RequiredPlayer: 4,
	})
	games = append(games, Game{
		ID: "a5",
		Title: "",
		RequiredPlayer: 5,
	})

}


func main() {
	createGames()
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/connect", playerConnection)
	router.GET("/players", getPlayers)
	router.GET("/games", getGames)
	router.GET("/queue", getQueue)
	router.POST("/queue/join", enterQueue)

	if err := router.Run(":4000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}