package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// PlayerStatus enum
type PlayerStatus string

const (
	PLAYING PlayerStatus = "Playing"
	ONLINE  PlayerStatus = "Online"
	OFFLINE PlayerStatus = "Offline"
)

// Player struct
type Player struct {
	Username     string       `json:"username"`
	UserId       string       `json:"userId"`
	PlayerStatus PlayerStatus `json:"playerStatus"`
	Password     string       `json:"-"` 
}

// Requests
type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LogOutRequest struct {
	Username string `json:"username"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Responses
type PlayerResponse struct {
	Username     string `json:"username"`
	UserId       string `json:"userId"`
	PlayerStatus string `json:"playerStatus"`
}

// PlayerService struct
type PlayerService struct {
	Players map[string]*Player
	mu      sync.RWMutex
}

// Constructor
func NewPlayerService() *PlayerService {
	return &PlayerService{
		Players: make(map[string]*Player),
	}
}

// Login method
func (ps *PlayerService) Login(request LoginRequest) (*PlayerResponse, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	player, exists := ps.Players[request.Username]
	if !exists {
		return nil, errors.New("player not found")
	}

	hashedPassword := hashPassword(request.Password)
	if hashedPassword != player.Password {
		return nil, errors.New("invalid credentials")
	}

	player.PlayerStatus = ONLINE
	return mapToPlayerResponse(player), nil
}

// Logout method
func (ps *PlayerService) Logout(request LogOutRequest) (string, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	player, exists := ps.Players[request.Username]
	if !exists {
		return "", errors.New("player not found")
	}

	if player.PlayerStatus == OFFLINE {
		return "", errors.New("player already offline")
	}

	player.PlayerStatus = OFFLINE
	message := fmt.Sprintf("Player: %v is set to Offline", player.Username)
	return message, nil
}

// Signup method
func (ps *PlayerService) Signup(request SignupRequest) (*PlayerResponse, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, exists := ps.Players[request.Username]; exists {
		return nil, fmt.Errorf("player with username %v already exists", request.Username)
	}

	newPlayer := &Player{
		Username:     request.Username,
		UserId:       uuid.New().String(),
		PlayerStatus: OFFLINE,
		Password:     hashPassword(request.Password),
	}

	ps.Players[request.Username] = newPlayer
	return mapToPlayerResponse(newPlayer), nil
}

// mapToPlayerResponse helper
func mapToPlayerResponse(player *Player) *PlayerResponse {
	return &PlayerResponse{
		Username:     player.Username,
		UserId:       player.UserId,
		PlayerStatus: string(player.PlayerStatus),
	}
}

// hashPassword helper
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
