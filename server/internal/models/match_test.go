package models

import (
	"testing"
	"time"
)

func TestNewMatch(t *testing.T) {
	gameID := "game-123"
	players := []string{"player1", "player2", "player3"}

	match := NewMatch(gameID, players)

	if match == nil {
		t.Fatal("Expected non-nil match")
	}
	if match.ID == "" {
		t.Error("Expected non-empty match ID")
	}
	if match.GameID != gameID {
		t.Errorf("Expected GameID %s, got %s", gameID, match.GameID)
	}
	if len(match.Players) != len(players) {
		t.Errorf("Expected %d players, got %d", len(players), len(match.Players))
	}
	for i, player := range players {
		if match.Players[i] != player {
			t.Errorf("Expected player %s at index %d, got %s", player, i, match.Players[i])
		}
	}
	if match.MatchState != MatchStateCreated {
		t.Errorf("Expected initial state %s, got %s", MatchStateCreated, match.MatchState)
	}
	if match.Duration != 0 {
		t.Errorf("Expected initial duration 0, got %f", match.Duration)
	}
	if time.Since(match.CreatedAt) > time.Second {
		t.Error("CreatedAt should be close to current time")
	}
}

func TestStartMatch(t *testing.T) {
	match := NewMatch("game-456", []string{"alice", "bob"})

	// Test starting from Created state
	match.StartMatch()
	if match.MatchState != MatchStateRunning {
		t.Errorf("Expected state %s after start, got %s", MatchStateRunning, match.MatchState)
	}

	// Test starting from Running state (should remain Running)
	match.StartMatch()
	if match.MatchState != MatchStateRunning {
		t.Errorf("Expected state to remain %s, got %s", MatchStateRunning, match.MatchState)
	}

	// Test starting from Ended state (should not change)
	match.MatchState = MatchStateEnded
	match.StartMatch()
	if match.MatchState != MatchStateEnded {
		t.Errorf("Expected state to remain %s, got %s", MatchStateEnded, match.MatchState)
	}
}

func TestEndMatch(t *testing.T) {
	match := NewMatch("game-789", []string{"player1", "player2"})

	// Test ending from Created state (should not change)
	match.EndMatch()
	if match.MatchState != MatchStateCreated {
		t.Errorf("Expected state to remain %s, got %s", MatchStateCreated, match.MatchState)
	}
	if match.Duration != 0 {
		t.Errorf("Expected duration to remain 0, got %f", match.Duration)
	}

	// Test ending from Running state
	match.MatchState = MatchStateRunning
	time.Sleep(10 * time.Millisecond) // Ensure some time passes
	match.EndMatch()

	if match.MatchState != MatchStateEnded {
		t.Errorf("Expected state %s after end, got %s", MatchStateEnded, match.MatchState)
	}
	if match.Duration <= 0 {
		t.Errorf("Expected positive duration, got %f", match.Duration)
	}

	// Test ending from already Ended state (should not change duration)
	previousDuration := match.Duration
	time.Sleep(5 * time.Millisecond)
	match.EndMatch()
	if match.Duration != previousDuration {
		t.Errorf("Expected duration to remain %f, got %f", previousDuration, match.Duration)
	}
}

func TestSetWinner(t *testing.T) {
	match := NewMatch("game-winner", []string{"player1", "player2"})

	// Test setting valid winners
	match.SetWinner("BlackTeam")
	if match.MatchInfo.Winner != "BlackTeam" {
		t.Errorf("Expected winner BlackTeam, got %s", match.MatchInfo.Winner)
	}

	match.SetWinner("WhiteTeam")
	if match.MatchInfo.Winner != "WhiteTeam" {
		t.Errorf("Expected winner WhiteTeam, got %s", match.MatchInfo.Winner)
	}

	// Test setting invalid winner (should not change)
	match.SetWinner("InvalidTeam")
	if match.MatchInfo.Winner != "WhiteTeam" {
		t.Errorf("Expected winner to remain WhiteTeam, got %s", match.MatchInfo.Winner)
	}

	match.SetWinner("")
	if match.MatchInfo.Winner != "WhiteTeam" {
		t.Errorf("Expected winner to remain WhiteTeam, got %s", match.MatchInfo.Winner)
	}
}

func TestIsActive(t *testing.T) {
	match := NewMatch("game-active", []string{"player1", "player2"})

	// Test Created state
	if match.IsActive() {
		t.Error("Expected match to not be active in Created state")
	}

	// Test Running state
	match.MatchState = MatchStateRunning
	if !match.IsActive() {
		t.Error("Expected match to be active in Running state")
	}

	// Test Ended state
	match.MatchState = MatchStateEnded
	if match.IsActive() {
		t.Error("Expected match to not be active in Ended state")
	}
}

func TestIsCompleted(t *testing.T) {
	match := NewMatch("game-completed", []string{"player1", "player2"})

	// Test Created state
	if match.IsCompleted() {
		t.Error("Expected match to not be completed in Created state")
	}

	// Test Running state
	match.MatchState = MatchStateRunning
	if match.IsCompleted() {
		t.Error("Expected match to not be completed in Running state")
	}

	// Test Ended state
	match.MatchState = MatchStateEnded
	if !match.IsCompleted() {
		t.Error("Expected match to be completed in Ended state")
	}
}

func TestMatchLifecycle(t *testing.T) {
	// Test complete match lifecycle
	match := NewMatch("game-lifecycle", []string{"alice", "bob"})
	
	// Initial state
	if match.MatchState != MatchStateCreated {
		t.Fatalf("Expected initial state %s, got %s", MatchStateCreated, match.MatchState)
	}
	if match.IsActive() {
		t.Error("Match should not be active initially")
	}
	if match.IsCompleted() {
		t.Error("Match should not be completed initially")
	}

	// Start match
	match.StartMatch()
	if match.MatchState != MatchStateRunning {
		t.Fatalf("Expected state %s after start, got %s", MatchStateRunning, match.MatchState)
	}
	if !match.IsActive() {
		t.Error("Match should be active after start")
	}
	if match.IsCompleted() {
		t.Error("Match should not be completed while running")
	}

	// Set up teams and moves
	match.MatchInfo.BlackTeam = []string{"alice"}
	match.MatchInfo.WhiteTeam = []string{"bob"}
	match.MatchInfo.MoveInfo = "e4 e5 Nf3"

	// Set winner
	match.SetWinner("BlackTeam")
	if match.MatchInfo.Winner != "BlackTeam" {
		t.Errorf("Expected winner BlackTeam, got %s", match.MatchInfo.Winner)
	}

	// End match
	time.Sleep(5 * time.Millisecond) // Ensure some duration
	match.EndMatch()
	
	if match.MatchState != MatchStateEnded {
		t.Fatalf("Expected state %s after end, got %s", MatchStateEnded, match.MatchState)
	}
	if match.IsActive() {
		t.Error("Match should not be active after end")
	}
	if !match.IsCompleted() {
		t.Error("Match should be completed after end")
	}
	if match.Duration <= 0 {
		t.Errorf("Expected positive duration, got %f", match.Duration)
	}
}

func TestMatchInfoFields(t *testing.T) {
	match := NewMatch("game-info", []string{"player1", "player2", "player3", "player4"})

	// Test setting team information
	blackTeam := []string{"player1", "player2"}
	whiteTeam := []string{"player3", "player4"}
	moveInfo := "1. e4 e5 2. Nf3 Nc6 3. Bb5"

	match.MatchInfo.BlackTeam = blackTeam
	match.MatchInfo.WhiteTeam = whiteTeam
	match.MatchInfo.MoveInfo = moveInfo

	if len(match.MatchInfo.BlackTeam) != len(blackTeam) {
		t.Errorf("Expected %d black team members, got %d", len(blackTeam), len(match.MatchInfo.BlackTeam))
	}
	if len(match.MatchInfo.WhiteTeam) != len(whiteTeam) {
		t.Errorf("Expected %d white team members, got %d", len(whiteTeam), len(match.MatchInfo.WhiteTeam))
	}
	if match.MatchInfo.MoveInfo != moveInfo {
		t.Errorf("Expected move info '%s', got '%s'", moveInfo, match.MatchInfo.MoveInfo)
	}

	// Verify team members
	for i, player := range blackTeam {
		if match.MatchInfo.BlackTeam[i] != player {
			t.Errorf("Expected black team player %s at index %d, got %s", player, i, match.MatchInfo.BlackTeam[i])
		}
	}
	for i, player := range whiteTeam {
		if match.MatchInfo.WhiteTeam[i] != player {
			t.Errorf("Expected white team player %s at index %d, got %s", player, i, match.MatchInfo.WhiteTeam[i])
		}
	}
}

func TestMatchConstants(t *testing.T) {
	// Test that match state constants are correct
	if MatchStateCreated != "Created" {
		t.Errorf("Expected MatchStateCreated to be 'Created', got %s", MatchStateCreated)
	}
	if MatchStateRunning != "Running" {
		t.Errorf("Expected MatchStateRunning to be 'Running', got %s", MatchStateRunning)
	}
	if MatchStateEnded != "Ended" {
		t.Errorf("Expected MatchStateEnded to be 'Ended', got %s", MatchStateEnded)
	}
}

func TestConcurrentMatches(t *testing.T) {
	// Test that multiple matches can be created independently
	match1 := NewMatch("game-1", []string{"alice", "bob"})
	match2 := NewMatch("game-2", []string{"charlie", "dave"})

	if match1.ID == match2.ID {
		t.Error("Expected different IDs for different matches")
	}
	if match1.GameID == match2.GameID {
		t.Error("Expected different GameIDs")
	}
	
	// Test independent state changes
	match1.StartMatch()
	if match2.IsActive() {
		t.Error("Match2 should not be affected by match1 state change")
	}
	
	match2.StartMatch()
	match1.EndMatch()
	if !match2.IsActive() {
		t.Error("Match2 should still be active after match1 ends")
	}
}