package api

import (
	"game-server/internal/handler"
	"game-server/internal/service"

	"github.com/gorilla/mux"
)

func PlayerRegisterRoutes(router *mux.Router) {
	// create service
	playerService := service.NewPlayerService()

	// inject service into handler
	playerHandler := handler.NewPlayerHandler(playerService)

	// register routes
	router.HandleFunc("/api/login", playerHandler.Login).Methods("POST")
	router.HandleFunc("/api/logout", playerHandler.Logout).Methods("POST")
	router.HandleFunc("/api/signup", playerHandler.SignUp).Methods("POST")
}
