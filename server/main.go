package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Player struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Level float32 `json:"level"`
}

type Game struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	RequiredPlayer int    `json:"requiredPlayer"`
}

type ChatMessage struct {
	MatchID string `json:"matchId"`
	PlayerID string `json:"playerId"`
	Content string `json:"content"`
	Timestamp string `json:"timestamp"`
}

var (
	players []Player
	games   []Game
	queue   map[string][]string
	allMatch	map[string][]string
	runningMatch map[string][]string
	mutex sync.Mutex

	connections map[string]map[string]*websocket.Conn
	connMutex sync.Mutex
	upgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool{
			return true
		},
	}


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

func getGames(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"games": games,
	})
}

func getQueue(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"queue": queue,
	})
}

func enterQueue(c *gin.Context) {
	var req struct {
		PlayerID string `json:"playerId"`
		GameID   string `json:"gameId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var player Player
	playerFound := false

	for _, p := range players {
		if p.ID == req.PlayerID {
			player = p
			playerFound = true
			break
		}
	}

	if !playerFound {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("player with ID '%s' not found", req.PlayerID)})
		return
	}

	var game Game
	gameFound := false

	for _, g := range games {
		if g.ID == req.GameID {
			game = g
			gameFound = true
			break
		}
	}

	if !gameFound {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("game with ID '%s' not found", req.GameID)})
		return
	}

	queue[game.ID] = append(queue[game.ID], player.ID)

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("player(%s) join the queue for game %s(%s)", req.PlayerID, game.Title, req.GameID)})
}

func createGames() {
	games = append(games, Game{
		ID:             "a2",
		Title:          "Tic Tac Toe",
		RequiredPlayer: 2,
	})
	games = append(games, Game{
		ID:             "a4",
		Title:          "Ludu",
		RequiredPlayer: 4,
	})
	games = append(games, Game{
		ID:             "a5",
		Title:          "minecraft",
		RequiredPlayer: 5,
	})
}


func getGameForUser( c* gin.Context){
	userID := c.Param("userId")

	playerMatchID := ""
	for matchID, players := range runningMatch {
		for _, playerID := range  players {
			if playerID == userID {
				playerMatchID = matchID
			}
		}
	}
	if playerMatchID == ""{
		c.JSON(http.StatusNotFound, gin.H{
			"messages": "player are not in running match",
		})
		return 
	}
	c.JSON(http.StatusOK, gin.H{
		"matchID": playerMatchID,
	})
	

}

func getMatch(c *gin.Context){
	matchID := c.Param("matchId")
	playerID := c.Param("userId")

	players, exists := runningMatch[matchID]
	if !exists{
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("match %s not found", matchID),
		})
	}

	userInMatch := false
	for _, gamePlayerID := range players{
		if playerID == gamePlayerID {
			userInMatch = true
			break
		}
	}

	if !userInMatch{
		c.JSON(http.StatusForbidden, gin.H{
			"message": fmt.Sprintf("player %s is not in match %s", playerID, matchID),
		})
	}



	//  send the client all necessary resources to the client
	c.JSON(http.StatusOK, gin.H{
		"matchID": matchID,
		"players": players,
	})


}


func matchMake(){
	for {
		mutex.Lock();
		for _, game := range games{
			if queuedPlayers, exists := queue[game.ID]; exists && len(queuedPlayers) >= game.RequiredPlayer {
				matchID := fmt.Sprintf("match-%d", time.Now().UnixNano())
				newMatch := queuedPlayers[:game.RequiredPlayer]
				queue[game.ID] = queuedPlayers[game.RequiredPlayer:]

				runningMatch[matchID] = newMatch
				allMatch[matchID] = newMatch
				log.Printf("Created match %s for game %s with players: %v", matchID, game.ID, newMatch)
			}
		}
		mutex.Unlock()
		time.Sleep(1 * time.Second)

	}
}

func handleChat(c *gin.Context){
	matchID := c.Param("matchId")
	playerID := c.Param("playerId")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil{
		log.Printf("websocket upgrader error: %v\n", err)
		return 
	}

	connMutex.Lock()
	if _, exists := connections[matchID]; !exists{
		connections[matchID] = make(map[string]*websocket.Conn)
	}
	connections[matchID][playerID] = conn
	connMutex.Unlock()

	log.Printf("Player %s connected to chat for match %s", playerID, matchID)

	for{
		var msg ChatMessage
		err := conn.ReadJSON(&msg)
		if err != nil{
			log.Printf("websocket read error for player %s: %v", playerID, err)
			break;
		}
		if msg.MatchID != matchID || msg.PlayerID != playerID {
			continue 

		}

		msg.Timestamp = time.Now().Format(time.RFC3339)
		connMutex.Lock()
		for pID, pConn := range connections[matchID] {
			if pID != playerID {
				if err := pConn.WriteJSON(msg); err != nil {
					log.Printf("Error sending message to player %s: %v", pID, err)
				}
			}
		}
		connMutex.Unlock()
	}
	connMutex.Lock()
	delete(connections[matchID],playerID)
	if len(connections[matchID]) == 0 {
		delete(connections, matchID)
	}
	connMutex.Unlock()
	conn.Close()
	log.Printf("Player %s disconnected from chat for match %s.", playerID, matchID)

}

func init(){
	connections = make(map[string]map[string]*websocket.Conn)
}

func main() {
	queue = make(map[string][]string)
	runningMatch = make(map[string][]string)
	allMatch = make(map[string][]string)
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
	router.GET("/match/:userId", getGameForUser)
	router.GET("/running-match/:matchId/:userId", getMatch)
	router.GET("/chat/:matchId/:userId", handleChat)


	go matchMake();

	if err := router.Run(":4000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
