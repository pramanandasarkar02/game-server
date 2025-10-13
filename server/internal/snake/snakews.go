package snake

import (
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
		handlePlayerInput(matchId, playerId, string(message))
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


func boradcastToMatch(matchId string, message []byte){
	matchConnMutex.Lock()
	defer matchConnMutex.Unlock()

	for _, conn := range matchConnections[matchId] {
		conn.WriteMessage(websocket.TextMessage, message)
	}
}

func handlePlayerInput(matchId, playerId, input string){

	log.Println(matchId, playerId, input)

	newState := []byte(`{"type":"update","matchId":"` + matchId + `"}`)
	boradcastToMatch(matchId, newState)
}