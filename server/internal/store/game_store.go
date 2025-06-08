package store


import (
	"sync"

	"github.com/pramanandasarkar02/game-server/internal/game"
)

type GameStore struct {
	games map[string]game.Game
	mutex sync.RWMutex
}

func NewGameStore() *GameStore {
	return &GameStore{
		games: make(map[string]game.Game),
	}
}

func (s *GameStore) AddGame(g game.Game) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.games[g.ID()] = g
}

func (s *GameStore) GetGame(id string) (game.Game, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	g, exists := s.games[id]
	return g, exists
}

func (s *GameStore) GetAll() []game.Game {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	games := make([]game.Game, 0, len(s.games))
	for _, g := range s.games {
		games = append(games, g)
	}
	return games
}