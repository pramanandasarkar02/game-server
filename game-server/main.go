
package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gin-contrib/cors"
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
	MatchID   string `json:"matchId"`
	PlayerID  string `json:"playerId"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type TicTacToeState struct {
	Board     [9]string `json:"board"`
	Turn      string    `json:"turn"` // Player ID of whose turn it is
	Winner    string    `json:"winner"`
	IsDraw    bool      `json:"isDraw"`
	Players   []string  `json:"players"`
}

type WebSocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

var (
	players      []Player
	games        []Game
	queue        map[string][]string
	allMatch     map[string][]string
	runningMatch map[string][]string
	mutex        sync.RWMutex

	connections map[string]map[string]*websocket.Conn
	connMutex   sync.RWMutex
	gameStates  map[string]*TicTacToeState
	stateMutex  sync.RWMutex
	upgrader    = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
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
	
	mutex.Lock()
	players = append(players, newPlayer)
	mutex.Unlock()

	log.Printf("%s is connected", newPlayer.Name)

	c.JSON(http.StatusOK, gin.H{
		"message": "Player connected successfully",
		"player":  newPlayer,
	})
}

func getPlayers(c *gin.Context) {
	mutex.RLock()
	playersCopy := make([]Player, len(players))
	copy(playersCopy, players)
	mutex.RUnlock()
	
	c.JSON(http.StatusOK, gin.H{
		"players": playersCopy,
	})
}

func getGames(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"games": games,
	})
}

func getQueue(c *gin.Context) {
	mutex.RLock()
	queueCopy := make(map[string][]string)
	for k, v := range queue {
		queueCopy[k] = make([]string, len(v))
		copy(queueCopy[k], v)
	}
	mutex.RUnlock()
	
	c.JSON(http.StatusOK, gin.H{
		"queue": queueCopy,
	})
}

func enterQueue(c *gin.Context) {
	var req struct {
		PlayerID string `json:"playerID"`
		GameID   string `json:"gameID"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mutex.RLock()
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
		mutex.RUnlock()
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
	mutex.RUnlock()

	if !gameFound {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("game with ID '%s' not found", req.GameID)})
		return
	}

	mutex.Lock()
	queue[game.ID] = append(queue[game.ID], player.ID)
	mutex.Unlock()

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

func getGameForUser(c *gin.Context) {
	userID := c.Param("userId")

	mutex.RLock()
	playerMatchID := ""
	for matchID, players := range runningMatch {
		for _, playerID := range players {
			if playerID == userID {
				playerMatchID = matchID
				break
			}
		}
		if playerMatchID != "" {
			break
		}
	}
	mutex.RUnlock()

	if playerMatchID == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "player is not in running match",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"matchID": playerMatchID,
	})
}

func getMatch(c *gin.Context) {
	matchID := c.Param("matchId")
	playerID := c.Param("userId")

	mutex.RLock()
	players, exists := runningMatch[matchID]
	if !exists {
		mutex.RUnlock()
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("match %s not found", matchID),
		})
		return
	}

	userInMatch := false
	for _, gamePlayerID := range players {
		if playerID == gamePlayerID {
			userInMatch = true
			break
		}
	}
	
	playersCopy := make([]string, len(players))
	copy(playersCopy, players)
	mutex.RUnlock()

	if !userInMatch {
		c.JSON(http.StatusForbidden, gin.H{
			"message": fmt.Sprintf("player %s is not in match %s", playerID, matchID),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"matchID": matchID,
		"players": playersCopy,
	})
}

func matchMake() {
	for {
		mutex.Lock()
		for _, game := range games {
			if queuedPlayers, exists := queue[game.ID]; exists && len(queuedPlayers) >= game.RequiredPlayer {
				matchID := fmt.Sprintf("match-%d", time.Now().UnixNano())
				newMatch := make([]string, game.RequiredPlayer)
				copy(newMatch, queuedPlayers[:game.RequiredPlayer])
				queue[game.ID] = queuedPlayers[game.RequiredPlayer:]

				runningMatch[matchID] = newMatch
				allMatch[matchID] = newMatch
				log.Printf("Created match %s for game %s with players: %v", matchID, game.ID, newMatch)

				// Initialize Tic-Tac-Toe state for the match
				if game.ID == "a2" {
					stateMutex.Lock()
					gameStates[matchID] = &TicTacToeState{
						Board:   [9]string{},
						Turn:    newMatch[0], // First player starts
						Players: newMatch,
					}
					stateMutex.Unlock()
				}
			}
		}
		mutex.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func checkWin(board [9]string) (string, bool) {
	winPatterns := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // Rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // Columns
		{0, 4, 8}, {2, 4, 6}, // Diagonals
	}
	for _, pattern := range winPatterns {
		a, b, c := pattern[0], pattern[1], pattern[2]
		if board[a] != "" && board[a] == board[b] && board[a] == board[c] {
			return board[a], true
		}
	}
	if isDraw := !contains(board[:], ""); isDraw {
		return "", true
	}
	return "", false
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func handleChat(c *gin.Context) {
	matchID := c.Param("matchId")
	playerID := c.Param("playerId")

	// Validate player is in the match
	mutex.RLock()
	players, exists := runningMatch[matchID]
	mutex.RUnlock()
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("match %s not found", matchID)})
		return
	}
	
	userInMatch := false
	for _, pID := range players {
		if pID == playerID {
			userInMatch = true
			break
		}
	}
	if !userInMatch {
		c.JSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("player %s is not in match %s", playerID, matchID)})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error for player %s in match %s: %v", playerID, matchID, err)
		return
	}

	// Set up ping/pong handlers
	conn.SetPingHandler(func(appData string) error {
		log.Printf("Received ping from player %s in match %s", playerID, matchID)
		return conn.WriteMessage(websocket.PongMessage, []byte(appData))
	})

	conn.SetPongHandler(func(appData string) error {
		log.Printf("Received pong from player %s in match %s", playerID, matchID)
		return nil
	})

	// Set read deadline for connection health
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	connMutex.Lock()
	if _, exists := connections[matchID]; !exists {
		connections[matchID] = make(map[string]*websocket.Conn)
	}
	connections[matchID][playerID] = conn
	connMutex.Unlock()

	log.Printf("Player %s connected to chat for match %s", playerID, matchID)

	// Send initial game state
	stateMutex.RLock()
	if state, exists := gameStates[matchID]; exists {
		conn.WriteJSON(WebSocketMessage{
			Type: "state",
			Data: state,
		})
	}
	stateMutex.RUnlock()

	// Goroutine to send periodic pings
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				connMutex.RLock()
				if _, exists := connections[matchID][playerID]; !exists {
					connMutex.RUnlock()
					return
				}
				connMutex.RUnlock()
				
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("Error sending ping to player %s in match %s: %v", playerID, matchID, err)
					return
				}
			}
		}
	}()

	// Main message handling loop
	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("Player %s closed connection normally in match %s", playerID, matchID)
			} else {
				log.Printf("WebSocket read error for player %s in match %s: %v", playerID, matchID, err)
			}
			break
		}

		log.Printf("Received message from player %s in match %s: %+v", playerID, matchID, msg)

		switch msg.Type {
		case "chat":
			chatMsg, ok := msg.Data.(map[string]interface{})
			if !ok {
				log.Printf("Invalid chat message format from player %s in match %s", playerID, matchID)
				continue
			}
			content, ok := chatMsg["content"].(string)
			if !ok || content == "" {
				log.Printf("Empty or invalid chat message content from player %s in match %s", playerID, matchID)
				continue
			}

			chatMessage := ChatMessage{
				MatchID:   matchID,
				PlayerID:  playerID,
				Content:   content,
				Timestamp: time.Now().Format(time.RFC3339),
			}

			// Broadcast chat message
			connMutex.RLock()
			for pID, pConn := range connections[matchID] {
				if pID != playerID {
					if err := pConn.WriteJSON(WebSocketMessage{Type: "chat", Data: chatMessage}); err != nil {
						log.Printf("Error sending chat message to player %s in match %s: %v", pID, matchID, err)
					}
				}
			}
			connMutex.RUnlock()

		case "move":
			moveData, ok := msg.Data.(map[string]interface{})
			if !ok {
				log.Printf("Invalid move message format from player %s in match %s", playerID, matchID)
				continue
			}
			indexFloat, ok := moveData["index"].(float64)
			if !ok {
				log.Printf("Invalid move index from player %s in match %s", playerID, matchID)
				continue
			}
			index := int(indexFloat)

			stateMutex.Lock()
			state, exists := gameStates[matchID]
			if !exists || state.Winner != "" || state.IsDraw || state.Turn != playerID || index < 0 || index > 8 || state.Board[index] != "" {
				stateMutex.Unlock()
				continue
			}

			// Update game state
			symbol := "X"
			if state.Players[1] == playerID {
				symbol = "O"
			}
			state.Board[index] = symbol
			winner, gameOver := checkWin(state.Board)
			if gameOver {
				if winner != "" {
					state.Winner = playerID
				} else {
					state.IsDraw = true
				}
			} else {
				state.Turn = state.Players[0]
				if state.Turn == playerID {
					state.Turn = state.Players[1]
				}
			}

			// Broadcast updated state
			connMutex.RLock()
			for _, pConn := range connections[matchID] {
				if err := pConn.WriteJSON(WebSocketMessage{Type: "state", Data: state}); err != nil {
					log.Printf("Error sending state to player in match %s: %v", matchID, err)
				}
			}
			connMutex.RUnlock()
			stateMutex.Unlock()
		}
	}

	// Cleanup connection
	connMutex.Lock()
	if matchConnections, exists := connections[matchID]; exists {
		delete(matchConnections, playerID)
		if len(matchConnections) == 0 {
			delete(connections, matchID)
			stateMutex.Lock()
			delete(gameStates, matchID)
			stateMutex.Unlock()
		}
	}
	connMutex.Unlock()
	
	conn.Close()
	log.Printf("Player %s disconnected from chat for match %s", playerID, matchID)
}

func init() {
	connections = make(map[string]map[string]*websocket.Conn)
	gameStates = make(map[string]*TicTacToeState)
}

func main() {
	queue = make(map[string][]string)
	runningMatch = make(map[string][]string)
	allMatch = make(map[string][]string)
	createGames()
	
	router := gin.Default()

	router.Use(cors.Default())

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
	router.GET("/chat/:matchId/:playerId", handleChat)

	go matchMake()

	if err := router.Run(":4000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
