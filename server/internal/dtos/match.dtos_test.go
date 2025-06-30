package dtos

import (
	"testing"
	"time"
	"github.com/pramanandasarkar02/game-server/internal/models"
)

func TestToMatchResponse(t *testing.T) {
	// Create a test match
	match := &models.Match{
		ID:        "test-id",
		GameID:    "game-123",
		Players:   []string{"player1", "player2"},
		CreatedAt: time.Now(),
		Duration:  120.5,
		MatchInfo: models.MatchInfo{
			BlackTeam: []string{"player1"},
			WhiteTeam: []string{"player2"},
			Winner:    "BlackTeam",
			MoveInfo:  "e4 e5",
		},
		MatchState: models.MatchStateEnded,
	}

	response := ToMatchResponse(match)

	// Verify all fields are correctly mapped
	if response.ID != match.ID {
		t.Errorf("Expected ID %s, got %s", match.ID, response.ID)
	}
	if response.GameID != match.GameID {
		t.Errorf("Expected GameID %s, got %s", match.GameID, response.GameID)
	}
	if len(response.Players) != len(match.Players) {
		t.Errorf("Expected %d players, got %d", len(match.Players), len(response.Players))
	}
	if response.Duration != match.Duration {
		t.Errorf("Expected duration %f, got %f", match.Duration, response.Duration)
	}
	if response.MatchState != match.MatchState {
		t.Errorf("Expected match state %s, got %s", match.MatchState, response.MatchState)
	}
}

func TestFromCreateRequest(t *testing.T) {
	req := CreateMatchRequest{
		GameID:  "game-456",
		Players: []string{"alice", "bob", "charlie"},
	}

	match := FromCreateRequest(req)

	if match.GameID != req.GameID {
		t.Errorf("Expected GameID %s, got %s", req.GameID, match.GameID)
	}
	if len(match.Players) != len(req.Players) {
		t.Errorf("Expected %d players, got %d", len(req.Players), len(match.Players))
	}
	if match.MatchState != models.MatchStateCreated {
		t.Errorf("Expected match state %s, got %s", models.MatchStateCreated, match.MatchState)
	}
	if match.ID == "" {
		t.Error("Expected non-empty match ID")
	}
}

func TestApplyUpdate(t *testing.T) {
	match := models.NewMatch("game-789", []string{"player1", "player2"})

	updateReq := UpdateMatchRequest{
		MatchState: models.MatchStateRunning,
		Winner:     "BlackTeam",
		MoveInfo:   "d4 d5",
		BlackTeam:  []string{"player1"},
		WhiteTeam:  []string{"player2"},
	}

	ApplyUpdate(match, updateReq)

	if match.MatchState != models.MatchStateRunning {
		t.Errorf("Expected match state %s, got %s", models.MatchStateRunning, match.MatchState)
	}
	if match.MatchInfo.Winner != "BlackTeam" {
		t.Errorf("Expected winner BlackTeam, got %s", match.MatchInfo.Winner)
	}
	if match.MatchInfo.MoveInfo != "d4 d5" {
		t.Errorf("Expected move info 'd4 d5', got %s", match.MatchInfo.MoveInfo)
	}
	if len(match.MatchInfo.BlackTeam) != 1 || match.MatchInfo.BlackTeam[0] != "player1" {
		t.Errorf("Expected BlackTeam [player1], got %v", match.MatchInfo.BlackTeam)
	}
	if len(match.MatchInfo.WhiteTeam) != 1 || match.MatchInfo.WhiteTeam[0] != "player2" {
		t.Errorf("Expected WhiteTeam [player2], got %v", match.MatchInfo.WhiteTeam)
	}
}

func TestApplyUpdateWithEndState(t *testing.T) {
	match := models.NewMatch("game-end", []string{"player1", "player2"})

	// Simulate some time passing
	time.Sleep(10 * time.Millisecond)

	updateReq := UpdateMatchRequest{
		MatchState: models.MatchStateEnded,
		Winner:     "WhiteTeam",
	}

	ApplyUpdate(match, updateReq)

	if match.MatchState != models.MatchStateEnded {
		t.Errorf("Expected match state %s, got %s", models.MatchStateEnded, match.MatchState)
	}
	if match.Duration <= 0 {
		t.Errorf("Expected positive duration, got %f", match.Duration)
	}
	if match.MatchInfo.Winner != "WhiteTeam" {
		t.Errorf("Expected winner WhiteTeam, got %s", match.MatchInfo.Winner)
	}
}

func TestApplyUpdatePartialUpdate(t *testing.T) {
	match := models.NewMatch("game-partial", []string{"player1", "player2"})
	match.MatchInfo.MoveInfo = "initial moves"
	match.MatchInfo.BlackTeam = []string{"initial"}

	// Apply partial update - only MoveInfo
	updateReq := UpdateMatchRequest{
		MoveInfo: "updated moves",
	}

	ApplyUpdate(match, updateReq)

	if match.MatchInfo.MoveInfo != "updated moves" {
		t.Errorf("Expected move info 'updated moves', got %s", match.MatchInfo.MoveInfo)
	}
	// Verify other fields remain unchanged
	if len(match.MatchInfo.BlackTeam) != 1 || match.MatchInfo.BlackTeam[0] != "initial" {
		t.Errorf("Expected BlackTeam unchanged, got %v", match.MatchInfo.BlackTeam)
	}
	if match.MatchState != models.MatchStateCreated {
		t.Errorf("Expected match state unchanged %s, got %s", models.MatchStateCreated, match.MatchState)
	}
}

func TestMatchFilter(t *testing.T) {
	dateFrom := time.Now().Add(-24 * time.Hour)
	dateTo := time.Now()

	filter := MatchFilter{
		GameID:     "game-123",
		PlayerID:   "player1",
		MatchState: models.MatchStateRunning,
		DateFrom:   &dateFrom,
		DateTo:     &dateTo,
		Page:       1,
		Limit:      10,
	}

	// Test that filter fields are properly set
	if filter.GameID != "game-123" {
		t.Errorf("Expected GameID 'game-123', got %s", filter.GameID)
	}
	if filter.PlayerID != "player1" {
		t.Errorf("Expected PlayerID 'player1', got %s", filter.PlayerID)
	}
	if filter.MatchState != models.MatchStateRunning {
		t.Errorf("Expected MatchState %s, got %s", models.MatchStateRunning, filter.MatchState)
	}
	if filter.Page != 1 {
		t.Errorf("Expected Page 1, got %d", filter.Page)
	}
	if filter.Limit != 10 {
		t.Errorf("Expected Limit 10, got %d", filter.Limit)
	}
}