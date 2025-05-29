package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func createPlayer() {
	fmt.Print("Enter player ID: ")
	fmt.Scanln(&player.ID)
	fmt.Print("Enter Player Name: ")
	fmt.Scanln(&player.Name)
	fmt.Print("Enter player Level: ")
	fmt.Scanln(&player.Level)
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

func printCommand() {
	for {
		fmt.Println("\n================ Game Client ===============")
		fmt.Println("1. Create player")
		fmt.Println("2. Connect to server")
		fmt.Println("3. Check server liveness")
		fmt.Println("4. Exit")

		var query int
		fmt.Print("Enter Command: ")
		fmt.Scanln(&query)

		switch query {
		case 1:
			createPlayer()
		case 2:
			connectServer()
		case 3:
			checkServerLiveness()
		case 4:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid command. Please choose 1, 2, 3, or 4.")
		}
	}
}

func main() {
	printCommand()
}