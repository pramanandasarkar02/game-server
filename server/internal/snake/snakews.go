package snake

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var matchConnections = make(map[string]map[string]*websocket.Conn)
var matchConnMutex sync.Mutex

var snakeService = NewSnakeService()




func WsHandler(c *gin.Context) {

	matchId := c.Query("matchId")
	playerId := c.Query("playerId")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrading error:", err)
		return
	}
	defer conn.Close()

	registerConnection(matchId, playerId, conn)
	defer unregisterConnection(matchId, playerId)
	log.Printf("Player %s connected to match %s", playerId, matchId)
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


func boradcastToMatch(matchId string, chat PlayerChat){
	matchConnMutex.Lock()
	defer matchConnMutex.Unlock()

	for _, conn := range matchConnections[matchId] {
		conn.WriteMessage(websocket.TextMessage, []byte(chat.Message))
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

	boradcastToMatch(matchId, chat)
}