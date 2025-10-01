package handler

import (
	"encoding/json"
	"game-server/internal/service"
	"net/http"
)

type PlayerHandler struct {
	playerService *service.PlayerService
}

func NewPlayerHandler(ps *service.PlayerService) *PlayerHandler {
	return &PlayerHandler{
		playerService: ps,
	}
}

// Login handler
func (ph *PlayerHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req service.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	player, err := ph.playerService.Login(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(player)
}

// Logout handler
func (ph *PlayerHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req service.LogOutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	message, err := ph.playerService.Logout(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

// Signup handler
func (ph *PlayerHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req service.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	player, err := ph.playerService.Signup(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(player)
}
