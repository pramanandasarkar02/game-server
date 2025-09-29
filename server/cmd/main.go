package main

import (
	"game-server/internal/api"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	PORT := ":8080"

	api.PlayerRegisterRoutes(router)

	router.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Echo from game-server"))
	}).Methods("GET")

	log.Printf("Server Started at http://localhost%v", PORT)
	if err := http.ListenAndServe(PORT, router); err != nil {
		log.Println("Server failure: ", err)
	}
}
