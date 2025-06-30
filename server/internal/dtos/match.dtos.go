package dtos

import (
	"time"
	"github.com/pramanandasarkar02/game-server/internal/models"
)

// CreateMatchRequest represents the request payload for creating a new match
type CreateMatchRequest struct {
	GameID  string   `json:"gameId" validate:"required"`
	Players []string `json:"players" validate:"required,min=2"`
}

// UpdateMatchRequest represents the request payload for updating a match
type UpdateMatchRequest struct {
	MatchState models.MatchState `json:"matchState,omitempty"`
	Winner     string            `json:"winner,omitempty"`
	MoveInfo   string            `json:"moveInfo,omitempty"`
	BlackTeam  []string          `json:"blackTeam,omitempty"`
	WhiteTeam  []string          `json:"whiteTeam,omitempty"`
}

// MatchResponse represents the response payload for match operations
type MatchResponse struct {
	ID         string              `json:"id"`
	GameID     string              `json:"gameId"`
	Players    []string            `json:"players"`
	CreatedAt  time.Time           `json:"createdAt"`
	Duration   float64             `json:"duration"`
	MatchInfo  models.MatchInfo    `json:"matchInfo"`
	MatchState models.MatchState   `json:"matchState"`
}

// MatchListResponse represents the response for listing matches
type MatchListResponse struct {
	Matches []MatchResponse `json:"matches"`
	Total   int             `json:"total"`
	Page    int             `json:"page"`
	Limit   int             `json:"limit"`
}

// MatchFilter represents filters for querying matches
type MatchFilter struct {
	GameID     string              `json:"gameId,omitempty"`
	PlayerID   string              `json:"playerId,omitempty"`
	MatchState models.MatchState   `json:"matchState,omitempty"`
	DateFrom   *time.Time          `json:"dateFrom,omitempty"`
	DateTo     *time.Time          `json:"dateTo,omitempty"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
}

// ToMatchResponse converts a Match model to MatchResponse DTO
func ToMatchResponse(m *models.Match) MatchResponse {
	return MatchResponse{
		ID:         m.ID,
		GameID:     m.GameID,
		Players:    m.Players,
		CreatedAt:  m.CreatedAt,
		Duration:   m.Duration,
		MatchInfo:  m.MatchInfo,
		MatchState: m.MatchState,
	}
}

// FromCreateRequest creates a new Match from CreateMatchRequest
func FromCreateRequest(req CreateMatchRequest) *models.Match {
	return models.NewMatch(req.GameID, req.Players)
}

// ApplyUpdate applies the updates from UpdateMatchRequest to the match
func ApplyUpdate(m *models.Match, req UpdateMatchRequest) {
    // Store the original state before updating
    originalState := m.MatchState

    if req.MatchState != "" {
        m.MatchState = req.MatchState
    }

    if req.Winner != "" {
        m.SetWinner(req.Winner)
    }

    if req.MoveInfo != "" {
        m.MatchInfo.MoveInfo = req.MoveInfo
    }

    if len(req.BlackTeam) > 0 {
        m.MatchInfo.BlackTeam = req.BlackTeam
    }

    if len(req.WhiteTeam) > 0 {
        m.MatchInfo.WhiteTeam = req.WhiteTeam
    }

    // If match is being ended, calculate duration
    if req.MatchState == models.MatchStateEnded && originalState != models.MatchStateEnded {
        m.Duration = time.Since(m.CreatedAt).Seconds()
    }
}