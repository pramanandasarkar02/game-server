package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pramanandasarkar02/game-server/config"

	"github.com/pramanandasarkar02/game-server/internal/store"
	"github.com/pramanandasarkar02/game-server/pkg/logger"
)

type WebSocketManager struct {
	connections map[string]map[string]*websocket.Conn
	mutex       sync.RWMutex
	cfg         *config.Config
	gameStore   store.GameStore
	matchStore  store.MatchStore
}

type WebSocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type ChatMessage struct {
	MatchID   string `json:"matchId"`
	PlayerID  string `json:"playerId"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebSocketManager(cfg *config.Config, gameStore store.GameStore, matchStore store.MatchStore) *WebSocketManager {
	return &WebSocketManager{
		connections: make(map[string]map[string]*websocket.Conn),
		cfg:         cfg,
		gameStore:   gameStore,
		matchStore:  matchStore,
	}
}

func (m *WebSocketManager) HandleConnection(c *gin.Context) {
	matchID := c.Param("matchId")
	playerID := c.Param("playerId")

	// Validate match and player
	match, exists := m.matchStore.GetMatch(matchID)
	if !exists {
		c.JSON(404, gin.H{"message": fmt.Sprintf("match %s not found", matchID)})
		return
	}
	userInMatch := false
	for _, pID := range match.Players {
		if pID == playerID {
			userInMatch = true
			break
		}
	}
	if !userInMatch {
		c.JSON(403, gin.H{"message": fmt.Sprintf("player %s is not in match %s", playerID, matchID)})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("WebSocket upgrade error for player %s in match %s: %v", playerID, matchID, err)
		return
	}

	// Set up ping/pong handlers
	conn.SetPingHandler(func(appData string) error {
		logger.Info("Received ping from player %s in match %s", playerID, matchID)
		return conn.WriteMessage(websocket.PongMessage, []byte(appData))
	})
	conn.SetPongHandler(func(appData string) error {
		logger.Info("Received pong from player %s in match %s", playerID, matchID)
		return nil
	})

	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(time.Duration(m.cfg.WSReadTimeout) * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(time.Duration(m.cfg.WSReadTimeout) * time.Second))
		return nil
	})

	// Register connection
	m.mutex.Lock()
	if _, exists := m.connections[matchID]; !exists {
		m.connections[matchID] = make(map[string]*websocket.Conn)
	}
	m.connections[matchID][playerID] = conn
	m.mutex.Unlock()

	logger.Info("Player %s connected to WebSocket for match %s", playerID, matchID)

	// Send initial game state
	g, exists := m.gameStore.GetGame(match.GameID)
	if exists {
		if state := g.GetState(matchID); state != nil {
			if err := conn.WriteJSON(WebSocketMessage{Type: "state", Data: state}); err != nil {
				logger.Error("Error sending initial state to player %s in match %s: %v", playerID, matchID, err)
			}
		}
	}

	// Periodic ping
	go func() {
		ticker := time.NewTicker(time.Duration(m.cfg.WSPingInterval) * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			m.mutex.RLock()
			if _, exists := m.connections[matchID][playerID]; !exists {
				m.mutex.RUnlock()
				return
			}
			m.mutex.RUnlock()
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Error("Error sending ping to player %s in match %s: %v", playerID, matchID, err)
				return
			}
		}
	}()

	// Handle messages
	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				logger.Info("Player %s closed connection in match %s", playerID, matchID)
			} else {
				logger.Error("WebSocket read error for player %s in match %s: %v", playerID, matchID, err)
			}
			break
		}

		switch msg.Type {
		case "chat":
			chatMsg, ok := msg.Data.(map[string]interface{})
			if !ok || chatMsg["content"] == nil || chatMsg["content"].(string) == "" {
				logger.Warn("Invalid chat message from player %s in match %s", playerID, matchID)
				continue
			}
			chatMessage := ChatMessage{
				MatchID:   matchID,
				PlayerID:  playerID,
				Content:   chatMsg["content"].(string),
				Timestamp: time.Now().Format(time.RFC3339),
			}
			m.broadcast(matchID, playerID, WebSocketMessage{Type: "chat", Data: chatMessage})

		case "move":
			moveMsg, ok := msg.Data.(map[string]interface{})
			if !ok || moveMsg["index"] == nil {
				logger.Warn("Invalid move message from player %s in match %s", playerID, matchID)
				continue
			}
			g, exists := m.gameStore.GetGame(match.GameID)
			if !exists {
				logger.Warn("Game %s not found for match %s", match.GameID, matchID)
				continue
			}
			if err := g.HandleMove(c.Request.Context(), matchID, playerID, moveMsg); err != nil {
				logger.Warn("Move error for player %s in match %s: %v", playerID, matchID, err)
				continue
			}
			if state := g.GetState(matchID); state != nil {
				m.broadcast(matchID, "", WebSocketMessage{Type: "state", Data: state})
			}
		case "ping":
			// Respond to client ping
			if err := conn.WriteMessage(websocket.PongMessage, nil); err != nil {
				logger.Error("Error sending pong to player %s in match %s: %v", playerID, matchID, err)
			}
		}
	}

	// Cleanup
	m.mutex.Lock()
	if matchConnections, exists := m.connections[matchID]; exists {
		delete(matchConnections, playerID)
		if len(matchConnections) == 0 {
			delete(m.connections, matchID)
		}
	}
	m.mutex.Unlock()
	conn.Close()
	logger.Info("Player %s disconnected from WebSocket for match %s", playerID, matchID)
}

func (m *WebSocketManager) broadcast(matchID, excludePlayerID string, msg WebSocketMessage) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for pID, conn := range m.connections[matchID] {
		if pID != excludePlayerID {
			if err := conn.WriteJSON(msg); err != nil {
				logger.Error("Error broadcasting to player %s in match %s: %v", pID, matchID, err)
			}
		}
	}
}
