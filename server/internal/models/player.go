package models

import (
	"fmt"
	"time"
	"github.com/google/uuid"
)
type PlayerState string

const (
	PlayerStateInGame PlayerState = "InGame"
	PlayerStateInQuery PlayerState = "InQuery"
	PlayerStateOffline PlayerState = "Offline"
	PlayerStateOnline PlayerState = "Online"
)


func (ps PlayerState) IsValid() bool {
	switch ps {
	case PlayerStateInGame, PlayerStateInQuery, PlayerStateOffline, PlayerStateOnline:
		return true
	}
	return false
}

func (ps PlayerState) String() string {
	return string(ps)
}


type Player struct {
	ID           string            `json:"id"`
	Username     string            `json:"name"`
	Password     string            `json:"-"`
	Level        float64           `json:"level"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
	MatchHistory map[string]bool   `json:"matchHistory"`
	Score        int               `json:"score"`
	State        PlayerState       `json:"state"`
}


func NewPlayer(username string, password string) *Player {
	return &Player{
		ID:           GenerateUUID(),
		Username:     username,
		Password:     password,
		Level:        1.0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		MatchHistory: make(map[string]bool),
		Score:        0,
		State:        PlayerStateOffline,
	}
}

func (p *Player) UpdateLevel(newLevel float64) {
	p.Level = newLevel
	p.UpdatedAt = time.Now()
}

func (p *Player) AddMatch(matchID string, won bool) {
	if p.MatchHistory == nil {
		p.MatchHistory = make(map[string]bool)
	}
	p.MatchHistory[matchID] = won
	p.UpdatedAt = time.Now()
}

func (p *Player) SetState(state PlayerState) error {
	if !state.IsValid() {
		return fmt.Errorf("invalid player state: %s", state)
	}
	p.State = state
	p.UpdatedAt = time.Now()
	return nil
}


func (p *Player) GetMatchCount() int {
	return len(p.MatchHistory)
}
func (p *Player) UpdatePassword(newPassword string) {
	p.Password = newPassword
	p.UpdatedAt = time.Now()
}


func(p *Player) IsOnline() bool {
	return p.State == PlayerStateOnline || p.State == PlayerStateInGame || p.State == PlayerStateInQuery
}


// ----------------------------------------------- ********* ------------------------------------------------

func GenerateUUID() string {
	return uuid.New().String()
}