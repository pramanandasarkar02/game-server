package api

import (
	"game-server/internal/handler"
	"github.com/gorilla/mux"
)

func PlayerRegisterRoutes(router *mux.Router) {
	playerHandler := handler.NewPlayerHandler()
	
	router.HandleFunc("/api/login", playerHandler.Login).Methods("POST")
	router.HandleFunc("/api/logout", playerHandler.Logout).Methods("POST")
	router.HandleFunc("/api/signup", playerHandler.SignUp).Methods("POST")
}
