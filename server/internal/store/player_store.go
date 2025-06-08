package store

import (
	"sync"

	"github.com/pramanandasarkar02/game-server/internal/models"
)

type PlayerStore struct {
	players []models.Player
	mutex   sync.RWMutex
}

func NewPlayerStore() *PlayerStore {
	return &PlayerStore{
		players: make([]models.Player, 0),
	}
}

func (s *PlayerStore) AddPlayer(player models.Player) error {
	if err := player.Validate(); err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.players = append(s.players, player)
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
	return models.Player{}, false
}

func (s *PlayerStore) GetAll() []models.Player {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	playersCopy := make([]models.Player, len(s.players))
	copy(playersCopy, s.players)
	return playersCopy
}
