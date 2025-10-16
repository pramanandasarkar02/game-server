package snake

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	matchConnMutex sync.Mutex
	activeMatches = make(map[string]bool)
	activeMatchLock sync.Mutex
	snakeService = NewSnakeService()
	deadSnake = make([]string, 0)
)





func WsHandler(c *gin.Context) {
	playerId := c.Query("playerId")
	matchId := c.Query("matchId")
	
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrading error:", err)
		return
	}
	defer conn.Close()

	registerConnection(matchId, playerId, conn)
	defer unregisterConnection(matchId, playerId)
	log.Printf("Player %s connected to match %s", playerId, matchId)

	startMatchLoopOnce(matchId, playerId)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}


		// log.Printf("Received: %s\n", message)

		// if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		// 	log.Println("Error writing message:", err)
		// 	break
		// }
		handlePlayerInput(matchId, playerId, message)
	}
}


func registerConnection(matchId, playerId string, conn * websocket.Conn){
	matchConnMutex.Lock()
	defer matchConnMutex.Unlock()

	if matchConnections[matchId] == nil{
		matchConnections[matchId] = make(map[string]*websocket.Conn)
	}

	matchConnections[matchId][playerId] = conn
}

func unregisterConnection(matchId, playerId string) {
	matchConnMutex.Lock()
	defer matchConnMutex.Unlock()

	if matchConnections[matchId] != nil{
		delete(matchConnections[matchId], playerId)
	}
}


func broadcastChatToMatch(matchId string, chat PlayerChat){
	matchConnMutex.Lock()
	defer matchConnMutex.Unlock()

	for _, conn := range matchConnections[matchId] {
		conn.WriteMessage(websocket.TextMessage, []byte(chat.Message))
	}
}
func broadcastToMatch(matchId string, message []byte){
	matchConnMutex.Lock()
	defer matchConnMutex.Unlock()

	for _, conn := range matchConnections[matchId] {
		conn.WriteMessage(websocket.TextMessage, message)
	}
}

type PlayerMessage struct{
	Type string `json:"type"`
}

type PlayerMove struct{
	Type string `json:"type"`
	Direction Direction `json:"direction"`
}

type PlayerChat struct{
	Type string `json:"type"`
	From string `json:"from"`
	Message string	`json:"message"`
}



func handlePlayerInput(matchId, playerId string, input []byte){
	var msg PlayerMessage 
	if err := json.Unmarshal(input, &msg); err != nil{
		log.Printf("Invalid json from %s: %s", playerId, input)
	}

	switch msg.Type {
	case "move":
		handleMove(matchId, playerId, input)
	case "chat":
		handleChat(matchId, playerId, input)
	default:
		log.Printf("Unknown message type from %s: %s", playerId, msg)
	}
	


	// log.Println(matchId, playerId, input)

	// newState := []byte(`{"type":"update","matchId":"` + matchId + `"}`)
	// boradcastToMatch(matchId, newState)
}


func handleMove(matchId, playerId string, input []byte){
	var move PlayerMove
	if err := json.Unmarshal(input, &move); err != nil{
		log.Printf("Invalid move message from %s: %s", playerId, string(input))
		return 
	}
	snakeService.ExecuteMovement(matchId, playerId, move.Direction)
}
func handleChat(matchId, playerId string, input []byte){
	var chat PlayerChat
	if err := json.Unmarshal(input, &chat); err != nil{
		log.Printf("invalid chat message from %s: %s", playerId, string(input))
		return
	}

	broadcastChatToMatch(matchId, chat)
}

func startMatchLoopOnce(matchId, playerId string){
	activeMatchLock.Lock()
	defer activeMatchLock.Unlock()

	if activeMatches[matchId]{
		snakeService.AddPlayer(matchId, playerId)
		return 
	}
	snakeService.StartGame(matchId)
	snakeService.AddPlayer(matchId, playerId)

	activeMatches[matchId] = true
	go func() {
		defer func(){
			activeMatchLock.Lock()
			delete(activeMatches, matchId)
			activeMatchLock.Unlock()
			log.Printf("Match loop ended for %s", matchId)
		}()
		log.Printf("Starting match loop for %s", matchId)
		ticker100ms := time.NewTicker(1000 * time.Millisecond)
		ticker500ms := time.NewTicker(2000 * time.Millisecond)
		ticker1s := time.NewTicker(10 * time.Second)

		defer ticker100ms.Stop()
		defer ticker1s.Stop()
		defer ticker500ms.Stop()

		for {
			select {
			case <- ticker100ms.C:
				broadcastBoardState(matchId)
			case <- ticker1s.C:
				snakeService.GenerateFood(matchId)
			case <- ticker500ms.C:
				for _, id := range deadSnake {
					if id == playerId {
						continue
					}
				}
				isCol, msg := snakeService.RunSnake(matchId, playerId)
				if isCol{
					log.Printf("player %v done by %v", playerId, msg)
					deadSnake = append(deadSnake, playerId)
				}
			}

			matchConnMutex.Lock()
			activePlayers := len(matchConnections[matchId])
			matchConnMutex.Unlock()

			if activePlayers == 0{
				log.Printf("No active player in match %s - stopping loop", matchId)
				return 
			}
		}

	}()
}

func broadcastBoardState(matchId string){
	matchConnMutex.Lock()
	playerIds := make([]string, 0, len(matchConnections[matchId]))
	for pId := range matchConnections[matchId]{
		playerIds = append(playerIds, pId)
	}
	matchConnMutex.Unlock()
	for _, playerId := range playerIds {
		boardState := snakeService.GetBoardStats(matchId, playerId)
		fmt.Printf("Board State: match(%v)-player(%v)\n%v\n", matchId, playerId, boardState)
		stateJSON, err := json.Marshal(boardState)
		if err != nil{
			log.Println("Error Marshalling board state: ", err)
			continue 
		}
		broadcastToMatch(matchId, stateJSON)
		
	}
}