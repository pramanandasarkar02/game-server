package models

import (
	"time"
	"github.com/google/uuid"
)

type MatchInfo struct {
	BlackTeam []string `json:"blackTeam"`
	WhiteTeam []string `json:"whiteTeam"`
	Winner    string   `json:"winner"` // "BlackTeam" or "WhiteTeam" or empty
	MoveInfo  string   `json:"moveInfo"`
}

type MatchState string

const (
	MatchStateCreated MatchState = "Created"
	MatchStateRunning MatchState = "Running"
	MatchStateEnded   MatchState = "Ended"
)

type Match struct {
	ID         string     `json:"id"`
	GameID     string     `json:"gameId"`
	Players    []string   `json:"players"`
	CreatedAt  time.Time  `json:"createdAt"`
	Duration   float64    `json:"duration"`
	MatchInfo  MatchInfo  `json:"matchInfo"`
	MatchState MatchState `json:"matchState"`
}

// NewMatch creates a new match with the given gameID and players
func NewMatch(gameID string, players []string) *Match {
	return &Match{
		ID:         uuid.New().String(),
		GameID:     gameID,
		Players:    players,
		CreatedAt:  time.Now(),
		Duration:   0,
		MatchInfo:  MatchInfo{},
		MatchState: MatchStateCreated,
	}
}

// StartMatch transitions the match to running state
func (m *Match) StartMatch() {
	if m.MatchState == MatchStateCreated {
		m.MatchState = MatchStateRunning
	}
}

// EndMatch transitions the match to ended state and calculates duration
func (m *Match) EndMatch() {
	if m.MatchState == MatchStateRunning {
		m.MatchState = MatchStateEnded
		m.Duration = time.Since(m.CreatedAt).Seconds()
	}
}

// SetWinner sets the winner of the match
func (m *Match) SetWinner(winner string) {
	if winner == "BlackTeam" || winner == "WhiteTeam" {
		m.MatchInfo.Winner = winner
	}
}

// IsActive returns true if the match is currently running
func (m *Match) IsActive() bool {
	return m.MatchState == MatchStateRunning
}

// IsCompleted returns true if the match has ended
func (m *Match) IsCompleted() bool {
	return m.MatchState == MatchStateEnded
}