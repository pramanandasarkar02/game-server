package models

import (
	"fmt"
	"time"
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
	Name         string            `json:"name"`
	Password     string            `json:"-"`
	Email        string            `json:"email"`
	Level        float32           `json:"level"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
	MatchHistory map[string]bool   `json:"matchHistory"` // matchID -> won/lose
	State        PlayerState       `json:"state"`
}


func NewPlayer(name, password, email string) *Player {
	now := time.Now()
	return &Player{
		Name: name,
		Password: password,
		Email: email,
		Level: 1.0,
		CreatedAt: now,
		UpdatedAt: now,
		MatchHistory: make(map[string]bool),
		State: PlayerStateOffline,
	}
}

func (p *Player) UpdateLevel(newLevel float32) {
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
		return fmt.Errorf("Invalid player state: %s", state)
	}
	p.State = state
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Player) GetWinRate() float32 {
	if len(p.MatchHistory) == 0 {
		return 0.0
	}
	wins := 0
	for _, won := range p.MatchHistory {
		if won {
			wins++
		}
	}
	return float32(wins) / float32(len(p.MatchHistory))
}

func (p *Player) GetMatchCount() int {
	return len(p.MatchHistory)
}
func (p *Player) UpdatePassword(newPassword string) {
	p.password = newPassword
	p.UpdatedAt = time.Now()
}

func (p *Player) UpdateEmail(newEmail string) {
	p.Email = newEmail
	p.UpdatedAt = time.Now()
}

func(p *Player) IsOnline() bool {
	return p.State == PlayerStateOnline || p.State == PlayerStateInGame || p.State == PlayerStateInQuery
}

func (p *Player) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("player name cannot be empty")
	}
	if len(p.Name) < 3 || len(p.Name) > 50 {
		return fmt.Errorf("player name must be between 3 and 50 characters")
	}
	if p.Email == "" {
		return fmt.Errorf("player email cannot be empty")
	}
	if p.Level < 0 {
		return fmt.Errorf("player level cannot be negative")
	}
	if !p.State.IsValid() {
		return fmt.Errorf("invalid player state: %s", p.State)
	}
	return nil
}

// MarshalJSON custom JSON marshaling
func (p *Player) MarshalJSON() ([]byte, error) {
	type Alias Player
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		Alias:     (*Alias)(p),
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
	})
}

// UnmarshalJSON custom JSON unmarshaling
func (p *Player) UnmarshalJSON(data []byte) error {
	type Alias Player
	aux := &struct {
		*Alias
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		Alias: (*Alias)(p),
	}
	
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	
	var err error
	p.CreatedAt, err = time.Parse(time.RFC3339, aux.CreatedAt)
	if err != nil {
		return err
	}
	
	p.UpdatedAt, err = time.Parse(time.RFC3339, aux.UpdatedAt)
	if err != nil {
		return err
	}
	
	return nil
}


