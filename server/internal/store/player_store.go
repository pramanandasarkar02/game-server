package store

import (
	"errors"
	"fmt"
	"github.com/pramanandasarkar02/game-server/internal/models"
	"log"
	"sync"
)

// in memory player store
type PlayerStore struct {
	players        []models.Player   // for storing players
	activeSessions map[string]string // for storing active sessions
	playerMutex    sync.RWMutex
	// sessionMutex  sync.RWMutex
}

func (ps *PlayerStore) GetPlayer(d string) (any, any) {
	panic("unimplemented")
}

// create new instance of player store
func NewPlayerStore() *PlayerStore {
	return &PlayerStore{
		players:        make([]models.Player, 0),
		activeSessions: make(map[string]string),
	}
}

// add new player to store
func (ps *PlayerStore) AddPlayer(player models.Player) error {
	ps.playerMutex.Lock()
	defer ps.playerMutex.Unlock()
	ps.players = append(ps.players, player)
	log.Printf("player %s added to store\n", player.Username)
	fmt.Printf("player %s added to store\n", player.Username)
	return nil
}

// get player id by username
func (ps *PlayerStore) GetPlayerIdByUsername(username string) (string, error) {
	ps.playerMutex.RLock()
	defer ps.playerMutex.RUnlock()
	for _, player := range ps.players {
		if player.Username == username {
			return player.ID, nil
		}
	}
	return "", errors.New("player not found")
}

// get player by player id
func (ps *PlayerStore) GetPlayerById(playerId string) (*models.Player, error) {
	ps.playerMutex.RLock()
	defer ps.playerMutex.RUnlock()
	for _, player := range ps.players {
		if player.ID == playerId {
			return &player, nil
		}
	}
	return nil, errors.New("player not found")
}

// update player
func (ps *PlayerStore) UpdatePlayer(player models.Player) error {
	ps.playerMutex.Lock()
	defer ps.playerMutex.Unlock()
	for i, p := range ps.players {
		if p.ID == player.ID {
			ps.players[i] = player
			log.Printf("player %s updated in store\n", player.Username)
			fmt.Printf("player %s updated in store\n", player.Username)
			return nil
		}
	}
	return fmt.Errorf("player with name %s not found", player.Username)
}

// delete player
func (ps *PlayerStore) DeletePlayer(playerId string) error {
	ps.playerMutex.Lock()
	defer ps.playerMutex.Unlock()
	for i, p := range ps.players {
		if p.ID == playerId {
			ps.players = append(ps.players[:i], ps.players[i+1:]...)
			log.Printf("player %s deleted from store\n", p.Username)
			fmt.Printf("player %s deleted from store\n", p.Username)
			return nil
		}
	}
	return fmt.Errorf("player with name %s not found", playerId)
}

// get all players
func (ps *PlayerStore) GetPlayers() ([]models.Player) {
	ps.playerMutex.RLock()
	defer ps.playerMutex.RUnlock()
	return ps.players
}
