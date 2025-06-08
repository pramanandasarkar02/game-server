
package store

import (
	"sync"
)

type QueueStore struct {
	queues map[string][]string
	mutex  sync.RWMutex
}

func NewQueueStore() *QueueStore {
	return &QueueStore{
		queues: make(map[string][]string),
	}
}

func (s *QueueStore) AddPlayer(gameID, playerID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.queues[gameID] = append(s.queues[gameID], playerID)
}

func (s *QueueStore) GetQueue(gameID string) []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	players := make([]string, len(s.queues[gameID]))
	copy(players, s.queues[gameID])
	return players
}

func (s *QueueStore) RemovePlayers(gameID string, count int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.queues[gameID]) >= count {
		s.queues[gameID] = s.queues[gameID][count:]
	}
}

func (s *QueueStore) GetAllQueues() map[string][]string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	queuesCopy := make(map[string][]string)
	for gameID, players := range s.queues {
		queuesCopy[gameID] = make([]string, len(players))
		copy(queuesCopy[gameID], players)
	}
	return queuesCopy
}
