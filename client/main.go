package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/gorilla/websocket"
)

type Player struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Level float32 `json:"level"`
}

type ChatMessage struct{
	MatchID string `json:"matchId"`
	PlayerID string `json:"playerId"`
	Content string `json:"content"`
	Timestamp string `json:"timestamp"`
}

func (player Player) String() string {
	return fmt.Sprintf("Player:\n\tID: %s\n\tName: %s\n\tLevel: %.2f", player.ID, player.Name, player.Level)
}

var (
	player   Player
	baseURL  = "http://localhost:4000"
	client   = &http.Client{}
	matchID string
	wsDialer = websocket.Dialer{}
)


func automaticStringGenerator(minLen, maxLen int) string{
	length := minLen + rand.IntN(maxLen - minLen)
	name := make([]byte, length)
	for i:= 0; i < length; i++ {
		name[i] = byte(97 + rand.IntN(25))
	}
	return string(name)
}

func createPlayer() {
	// fmt.Print("Enter player ID: ")
	// fmt.Scanln(&player.ID)
	// fmt.Print("Enter Player Name: ")
	// fmt.Scanln(&player.Name)
	// fmt.Print("Enter player Level: ")
	// fmt.Scanln(&player.Level)
	
	player.ID = automaticStringGenerator(4,8)
	player.Name = "Alice"
	player.Level = 0.0
	fmt.Println("Player created:", player)
}

func connectServer() {
	if player.Name == "" {
		fmt.Println("Create a player first (option 1).")
		return
	}

	requestPayload := struct {
		Name string `json:"name"`
		ID   string `json:"id"`
		Level float32 `json:"level"`
	}{
		Name:  player.Name,
		ID:    player.ID,
		Level: player.Level,
	}

	playerJSON, err := json.Marshal(requestPayload)
	if err != nil {
		log.Printf("Error marshalling request payload to JSON: %v\n", err)
		return
	}

	resp, err := client.Post(baseURL+"/connect", "application/json", bytes.NewBuffer(playerJSON))
	if err != nil {
		log.Printf("Error connecting to server: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading server response: %v\n", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Successfully connected to server:", string(body))
	} else {
		fmt.Println("Failed to connect to server:", string(body))
	}
}

func checkServerLiveness() {
	resp, err := client.Get(baseURL + "/ping")
	if err != nil {
		log.Printf("Error checking server liveness: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading server response: %v\n", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Server is alive:", string(body))
	} else {
		fmt.Println("Server check failed:", string(body))
	}
}

func joinQueue(){
	requestPayload := struct{
		PlayerID string `json:"playerID"`
		GameID string `json:"gameID"`
	}{
		PlayerID: player.ID,
		GameID: "a2",
	}
	queueJSON, err := json.Marshal(requestPayload)
	if err != nil{
		log.Printf("Error to marshalling request playload to JSON: %v\n", err)
		return 
	}

	response, err := client.Post(baseURL+"/queue/join", "application/json", bytes.NewBuffer(queueJSON))

	if err != nil {
		log.Printf("Error connecting to server: %v\n", err)
		return 
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading server response: %v\n", err)
		return 
	}
	if response.StatusCode == http.StatusOK{
		fmt.Println("Server replay: ", string(body))
	}else{
		fmt.Println("Something wrong: ", string(body))
	}


}

func getGame(){
	response, err := client.Get(baseURL + fmt.Sprintf("/match/%v", player.ID))
	if err != nil{
		fmt.Println("Error: ", err)
		return 
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading server response: %v\n", err)
		return 
	}
	


	if response.StatusCode == http.StatusOK {
		var resp struct{
			MatchID string `json:"matchID"`
		}
		err := json.Unmarshal(body, &resp)
		if err != nil {
			log.Printf("Error decoding JSON response: %v\n", err)
			return
		}

		matchID = resp.MatchID
		log.Println("Match ID:", matchID)
		return
	}
	log.Printf("Unexpected response status: %d\nResponse: %s", response.StatusCode, string(body))

}

func joinGame(){
	response, err := client.Get(baseURL + fmt.Sprintf("/running-match/%v/%v", matchID ,player.ID))
	if err != nil{
		fmt.Println("Error: ", err)
		return 
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading server response: %v\n", err)
		return 
	}
	if response.StatusCode == http.StatusNotFound{
		log.Println("Response: ", body)

	}


	if response.StatusCode == http.StatusOK {
		var resp struct{
			MatchID string `json:"matchID"`
			Players []string `json:"players"`
		}
		err := json.Unmarshal(body, &resp)
		if err != nil {
			log.Printf("Error decoding JSON response: %v\n", err)
			return
		}
		var players []string
		matchID = resp.MatchID
		players = resp.Players
		log.Println("Game ID:", matchID)
		log.Println("PlayersID: ", players)
		return
	}
	log.Printf("Unexpected response status: %d\nResponse: %s", response.StatusCode, string(body))
	
}

func startChat() {
    if matchID == "" || player.ID == "" {
        fmt.Println("Cannot start chat: missing match ID or player ID")
        return
    }

    wsURL := fmt.Sprintf("ws://localhost:4000/chat/%s/%s", matchID, player.ID)
    conn, _, err := wsDialer.Dial(wsURL, nil)
    if err != nil {
        fmt.Printf("Failed to connect to chat: %v\n", err)
        return
    }
    defer func() {
        conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Client exiting"))
        conn.Close()
        fmt.Println("Chat connection closed")
    }()

    fmt.Printf("Connected to chat for match %s as player %s\n", matchID, player.ID)
    fmt.Println("Type your messages (type 'exit' to quit chat):")

    go func() {
        for {
            var msg ChatMessage
            err := conn.ReadJSON(&msg)
            if err != nil {
                fmt.Printf("Chat disconnected: %v\n", err)
                return
            }
            if msg.PlayerID != player.ID {
                fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.PlayerID, msg.Content)
            }
        }
    }()

    go func() {
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()
        for range ticker.C {
            err := conn.WriteMessage(websocket.PingMessage, nil)
            if err != nil {
                fmt.Printf("Error sending ping: %v\n", err)
                return
            }
        }
    }()

    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        input := strings.TrimSpace(scanner.Text())
        if input == "exit" {
            fmt.Println("Exiting chat...")
            return
        }
        if input != "" {
            msg := ChatMessage{
                MatchID:   matchID,
                PlayerID:  player.ID,
                Content:   input,
                Timestamp: time.Now().Format(time.RFC3339), // Optional, server overrides this
            }
            log.Printf("Sending message: %+v", msg)
            err := conn.WriteJSON(msg)
            if err != nil {
                fmt.Printf("Error sending message: %v\n", err)
                return
            }
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Printf("Error reading input: %v\n", err)
    }
}



// func startChat(){
// 	if matchID == ""{
// 		fmt.Println("join a game first")
// 		return 
// 	}

// 	wsURL := fmt.Sprintf("ws://localhost:4000/chat/%s/%s", matchID, player.ID)
// 	conn, _, err := wsDialer.Dial(wsURL, nil)
// 	if err != nil{
// 		log.Println("Error connecting to websocket: %v\n", err)
// 		return 
// 	}
// 	defer func(){
// 		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Client exiting"))
// 		conn.Close()
// 		fmt.Println("Chat Connection Closed")
// 	}()

// 	fmt.Println("Connected to chat for match ", matchID)
// 	fmt.Println("Type your messages (type 'exit' to quit chat):")

// 	go func(){
// 		for {
// 			var msg ChatMessage
// 			err := conn.ReadJSON(&msg)
// 			if err != nil {
// 				log.Println("Websocket read err: %v\n", err)
// 				return 
// 			}
// 			if msg.PlayerID != player.ID {
// 				fmt.Printf("[%s] %s: %s", msg.Timestamp, msg.PlayerID, msg.Content)
// 			}
// 		}
// 	}()

// 	go func(){
// 		ticker := time.NewTicker(30 * time.Second)
// 		defer ticker.Stop()
// 		for range ticker.C{
// 			err := conn.WriteMessage(websocket.PingMessage, nil)
// 			if err != nil{
// 				fmt.Printf("Error sending ping: %v\n", err)
// 				return 
// 			}
// 		}
// 	}()

// 	scanner := bufio.NewScanner(os.Stdin)

// 	for scanner.Scan(){
		
// 		input := strings.TrimSpace(scanner.Text())
// 		if input == "exit" {
// 			fmt.Println("EXiting Chat...")
// 			return 
// 		}

// 		if input != ""{
// 			msg := ChatMessage{
// 				MatchID: matchID,
// 				PlayerID: player.ID,
// 				Content: input,
// 			}
// 			err := conn.WriteJSON(msg)
// 			if err != nil{
// 				log.Println("Error in sending message: %v\n", err)
// 				return 
// 			}
// 		}
// 	}
// 	if err := scanner.Err(); err != nil {
// 		fmt.Printf("Error reading input: %v\n", err)
// 	}
// }




func printCommand() {
	fmt.Println("\n================ Game Client ===============")
	fmt.Println("1. Check Server")
	fmt.Println("2. Connect to server")
	fmt.Println("3. Join queue")
	fmt.Println("4. getGame")
	fmt.Println("5. Join Game")
	fmt.Println("6. start Chat")
	fmt.Println("0. Exit")
	fmt.Println("==============================================")

}

func excuteCommands(){
	for {
		fmt.Println("==============================================\nRequest:")
		var query int
		fmt.Print("Enter Command: ")
		fmt.Scanln(&query)
		fmt.Println("-----------------------------------------------\nResponse")

		switch query {
		case 1:
			checkServerLiveness()
		case 2:
			connectServer()
		case 3:
			joinQueue()
		case 4:
			getGame()
		case 5:
			joinGame()
		case 6: 
			startChat()
		case 0:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid command. Please choose 0, 1, 2, 3,4, 5 or 6.")
		}
		fmt.Println("==============================================")
		fmt.Println()

		
	}
}

func main() {
	createPlayer();
	printCommand()
	excuteCommands()
}