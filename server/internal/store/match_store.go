
package store

import (
	"sync"

	"github.com/pramanandasarkar02/game-server/internal/models"
)

type MatchStore struct {
	matches map[string]models.Match
	mutex   sync.RWMutex
}

func NewMatchStore() *MatchStore {
	return &MatchStore{
		matches: make(map[string]models.Match),
	}
}

func (s *MatchStore) AddMatch(match models.Match) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.matches[match.ID] = match
}

func (s *MatchStore) GetMatch(id string) (models.Match, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	match, exists := s.matches[id]
	return match, exists
}

func (s *MatchStore) GetMatchForPlayer(playerID string) (models.Match, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, match := range s.matches {
		for _, pID := range match.Players {
			if pID == playerID {
				return match, true
			}
		}
	}
	return models.Match{}, false
}
