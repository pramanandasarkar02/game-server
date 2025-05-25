package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

var (
	LOBBY_SERVICE_PORT = ":8083"
)


func main(){
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis!")

	ctx := context.Background()
	http.HandleFunc("/lobbies/join", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return 
		}
		var req struct{
			UserID string `json:"user_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return 
		}
		

		fmt.Println(req.UserID)

		err := client.LPush(ctx, "match-making-queue", req.UserID).Err()
		if err != nil {
			http.Error(w, "Failed to join queue", http.StatusInternalServerError)
			return 
		}


		// simple match making queue if queue has 2+ player
		if queueLen, _ := client.LLen(ctx, "match-making-queue").Result(); queueLen >= 2{
			players, _ := client.LPopCount(ctx, "match-making-queue", 2).Result()
			
			json.NewEncoder(w).Encode(map[string]interface{}{"game-id": "game_id" + players[0], "players": players})
		} else {
			json.NewEncoder(w).Encode(map[string]string{"message": "waiting for match"})
		}
	})
	



	log.Printf("Starting server on http://localhost%s", LOBBY_SERVICE_PORT)
	log.Fatal(http.ListenAndServe(LOBBY_SERVICE_PORT, nil))
	

}