package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Player struct{
	ID string `json: "ID"`
	Name string `json: "Name"`
	Level float32 `json: "Level"`
}

func (player Player) String() string {
	return fmt.Sprint("player:\n\tID:\t%v\n\tName:\t%v\n\tLevel:\t%v", player.ID, player.Name, player.Level)
}


var(
	player Player
	port = ":4000"
)

func createPlayer(){
	fmt.Print("Enter player ID: ")
	fmt.Scanln(&player.ID)
	fmt.Print("Enter Player Name: ")
	fmt.Scanln(&player.Name)
	fmt.Print("Enter player Level: ")
	fmt.Scanln(&player.Level)

}

func connectServer(){
	if player.Name == ""{
		fmt.Println("create a player first (option 1).")
		return 
	}

	requestPayload := struct{
		Name string `json: "name"`
	}{
		Name: player.Name,
	}

	playerJSON, err := json.Marshal(requestPayload)

	if err != nil {
		log.Println("Error marshalling request playload to json: %v\n", err)
		return 
	}
	
	
	
}

func printCommand(){
	fmt.Println("================ Game Client ===============")
	fmt.Println("1. create player ")
	fmt.Println("2. connect to server")
	fmt.Println("3. check server liveness")
	
	var query int 
	for {
		fmt.Print("Enter Command: ")
		fmt.Scanln(&query)
		if query == 3{

		}


	}
}

func main(){
	

}