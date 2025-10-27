package snake

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	matchConnections = make(map[string]map[string]*websocket.Conn)
	matchConnMutex   sync.RWMutex
	activeMatches    = make(map[string]bool)
	activeMatchLock  sync.RWMutex
	snakeService     = NewSnakeService()
)

func WsHandler(c *gin.Context) {
	playerId := c.Query("playerId")
	matchId := c.Query("matchId")

	if playerId == "" || matchId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "playerId and matchId required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrading error:", err)
		return
	}
	defer conn.Close()

	registerConnection(matchId, playerId, conn)
	defer unregisterConnection(matchId, playerId)
	log.Printf("Player %s connected to match %s", playerId, matchId)

	db, err := sql.Open("sqlite3", "./matches.db")
	if err != nil {
		log.Printf("DB open error: %v", err)
		return
	}
	defer db.Close()

	playerIds, err := getPlayersByMatchId(db, matchId)
	if err != nil {
		log.Printf("Failed to load playerIds for match %s: %v", matchId, err)
		return
	}
	log.Printf("Players in match %s: %v", matchId, playerIds)

	// Initialize the game first, then add the player
	startMatchLoopOnce(matchId, playerIds)
	snakeService.AddPlayer(matchId, playerId)
	
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from player %s: %v", playerId, err)
			break
		}
		handlePlayerInput(matchId, playerId, message)
	}
}

func registerConnection(matchId, playerId string, conn *websocket.Conn) {
	matchConnMutex.Lock()
	defer matchConnMutex.Unlock()

	if matchConnections[matchId] == nil {
		matchConnections[matchId] = make(map[string]*websocket.Conn)
	}

	matchConnections[matchId][playerId] = conn
}

func unregisterConnection(matchId, playerId string) {
	matchConnMutex.Lock()
	defer matchConnMutex.Unlock()

	if matchConnections[matchId] != nil {
		delete(matchConnections[matchId], playerId)
		if len(matchConnections[matchId]) == 0 {
			delete(matchConnections, matchId)
		}
	}
}

func broadcastChatToMatch(matchId string, chat PlayerChat) {
	matchConnMutex.RLock()
	conns := make(map[string]*websocket.Conn)
	for k, v := range matchConnections[matchId] {
		conns[k] = v
	}
	matchConnMutex.RUnlock()

	msg, err := json.Marshal(chat)
	if err != nil {
		log.Println("Error marshalling chat:", err)
		return
	}

	for playerId, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Printf("Error sending chat to %s: %v", playerId, err)
		}
	}
}

func broadcastToMatch(matchId string, message []byte) {
	matchConnMutex.RLock()
	conns := make([]*websocket.Conn, 0, len(matchConnections[matchId]))
	for _, conn := range matchConnections[matchId] {
		conns = append(conns, conn)
	}
	matchConnMutex.RUnlock()

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error broadcasting to match %s: %v", matchId, err)
		}
	}
}

type PlayerMessage struct {
	Type string `json:"type"`
}

type PlayerMove struct {
	Type      string `json:"type"`
	Direction string `json:"direction"`
}

type PlayerChat struct {
	Type    string `json:"type"`
	From    string `json:"from"`
	Message string `json:"message"`
}

func handlePlayerInput(matchId, playerId string, input []byte) {
	var msg PlayerMessage
	if err := json.Unmarshal(input, &msg); err != nil {
		log.Printf("Invalid json from %s: %s", playerId, input)
		return
	}

	switch msg.Type {
	case "move":
		handleMove(matchId, playerId, input)
	case "chat":
		handleChat(matchId, playerId, input)
	default:
		log.Printf("Unknown message type from %s: %s", playerId, msg.Type)
	}
}

func handleMove(matchId, playerId string, input []byte) {
	var move PlayerMove
	if err := json.Unmarshal(input, &move); err != nil {
		log.Printf("Invalid move message from %s: %s", playerId, string(input))
		return
	}
	log.Printf("Move from %s in match %s: %s", playerId, matchId, move.Direction)
	snakeService.ExecuteMovement(matchId, playerId, strToDirection(move.Direction))
}

func strToDirection(dir string) Direction {
	switch strings.ToUpper(dir) {
	case "RIGHT":
		return RIGHT
	case "LEFT":
		return LEFT
	case "UP":
		return UP
	case "DOWN":
		return DOWN
	default:
		return ""
	}
}

func handleChat(matchId, playerId string, input []byte) {
	var chat PlayerChat
	if err := json.Unmarshal(input, &chat); err != nil {
		log.Printf("invalid chat message from %s: %s", playerId, string(input))
		return
	}

	// Ensure sender is set
	chat.From = playerId
	chat.Type = "chat"

	broadcastChatToMatch(matchId, chat)
}

func startMatchLoopOnce(matchId string, playerIds []string) {
	activeMatchLock.Lock()
	defer activeMatchLock.Unlock()

	// Check if match is already running
	if activeMatches[matchId] {
		log.Printf("Match loop already running for %s", matchId)
		return
	}

	snakeService.StartGame(matchId, playerIds)
	activeMatches[matchId] = true

	go func() {
		defer func() {
			activeMatchLock.Lock()
			delete(activeMatches, matchId)
			activeMatchLock.Unlock()
			log.Printf("Match loop ended for %s", matchId)
		}()
		
		log.Printf("Starting match loop for %s", matchId)
		ticker100ms := time.NewTicker(100 * time.Millisecond)
		ticker500ms := time.NewTicker(3000 * time.Millisecond)
		ticker1s := time.NewTicker(1 * time.Second)

		defer ticker100ms.Stop()
		defer ticker500ms.Stop()
		defer ticker1s.Stop()

		for {
			select {
			case <-ticker100ms.C:
				broadcastBoardState(matchId)
			case <-ticker1s.C:
				snakeService.GenerateFood(matchId)
			case <-ticker500ms.C:
				snakeService.RunAllSnake(matchId)
			}

			matchConnMutex.RLock()
			activePlayers := len(matchConnections[matchId])
			matchConnMutex.RUnlock()

			if activePlayers == 0 {
				log.Printf("No active player in match %s - stopping loop", matchId)
				return
			}
		}
	}()
}

func broadcastBoardState(matchId string) {
	matchConnMutex.RLock()
	playerIds := make([]string, 0, len(matchConnections[matchId]))
	for pId := range matchConnections[matchId] {
		playerIds = append(playerIds, pId)
	}
	matchConnMutex.RUnlock()

	if len(playerIds) == 0 {
		return
	}

	// Get board state once (assuming it's the same for all players)
	boardState := snakeService.GetBoardStats(matchId, playerIds[0])
	
	stateJSON, err := json.Marshal(map[string]interface{}{
		"type":  "update",
		"state": boardState,
	})
	if err != nil {
		log.Println("Error marshalling board state:", err)
		return
	}

	broadcastToMatch(matchId, stateJSON)
}

func getPlayersByMatchId(db *sql.DB, matchId string) ([]string, error) {
	row := db.QueryRow(`SELECT players FROM matches WHERE match_id = ?`, matchId)

	var playerList string
	if err := row.Scan(&playerList); err != nil {
		return nil, fmt.Errorf("failed to query players for match %s: %v", matchId, err)
	}

	playerIds := strings.Split(playerList, ",")
	for i := range playerIds {
		playerIds[i] = strings.TrimSpace(playerIds[i])
	}

	return playerIds, nil
}