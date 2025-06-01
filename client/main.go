package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"

)

type Player struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Level float32 `json:"level"`
}

func (player Player) String() string {
	return fmt.Sprintf("Player:\n\tID: %s\n\tName: %s\n\tLevel: %.2f", player.ID, player.Name, player.Level)
}

var (
	player   Player
	baseURL  = "http://localhost:4000"
	client   = &http.Client{}
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

}

func joinGame(){
	
}




func printCommand() {
	fmt.Println("\n================ Game Client ===============")
	fmt.Println("1. Check Server")
	fmt.Println("2. Connect to server")
	fmt.Println("3. Join queue")
	fmt.Println("4. getGame")
	fmt.Println("5. JoinGame")
	fmt.Println("6. ")
	fmt.Println("0. Exit")
	fmt.Println("==============================================")

}

func excuteCommands(){
	for {
		fmt.Println("-----------------------------------------------")
		var query int
		fmt.Print("Enter Command: ")
		fmt.Scanln(&query)

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
			
		case 0:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid command. Please choose 0, 1, 2, 3,4, 5 or 6.")
		}
		fmt.Println("-----------------------------------------------")
	}
}

func main() {
	createPlayer();
	printCommand()
	excuteCommands()
}