package store

import (
	"log"
	"sync"

	"github.com/pramanandasarkar02/game-server/internal/dtos"
	"github.com/pramanandasarkar02/game-server/internal/models"
)

// players data save in memory
type PlayerStore struct {
	players []models.Player
	mutex   sync.RWMutex
}



func NewPlayerStore() *PlayerStore {
	return &PlayerStore{
		players: make([]models.Player, 0),
	}
}


func (s *PlayerStore) CreatePlayer(playerDto dtos.CreatePlayerDto) {
	if err := playerDto.Validate(); err != nil {
		return 
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	// save in memory and postgress db


}

func (s* PlayerStore) GetPlayerProfileInfo(playerProfileInfoDto dtos.PlayerProfileInfoDto) {
	
}

func (s *PlayerStore) AddPlayer(player models.Player) error {
	if err := player.Validate(); err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.players = append(s.players, player)
	log.Printf("Player added: %+v, Total players: %d", player, len(s.players))
	return nil
}

func (s *PlayerStore) GetPlayer(id string) (models.Player, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, p := range s.players {
		if p.ID == id {
			return p, true
		}
	}
	log.Printf("Player with ID %s not found", id)
	return models.Player{}, false
}

func (s *PlayerStore) GetAll() []models.Player {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	playersCopy := make([]models.Player, len(s.players))
	copy(playersCopy, s.players)
	log.Printf("Returning %d players: %+v", len(playersCopy), playersCopy)
	return playersCopy // Return the copy to prevent external modification
}