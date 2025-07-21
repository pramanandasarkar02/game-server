package models

import (
	// "encoding/json"
	// "testing"
	// "time"
)

// func TestNewPlayer(t *testing.T) {
// 	player := NewPlayer("testuser", "password123", "test@example.com")
	
// 	if player.Name != "testuser" {
// 		t.Errorf("Expected name 'testuser', got %s", player.Name)
// 	}
// 	if player.Password != "password123" {
// 		t.Errorf("Expected password 'password123', got %s", player.Password)
// 	}
// 	if player.Email != "test@example.com" {
// 		t.Errorf("Expected email 'test@example.com', got %s", player.Email)
// 	}
// 	if player.Level != 1.0 {
// 		t.Errorf("Expected level 1.0, got %f", player.Level)
// 	}
// 	if player.State != PlayerStateOffline {
// 		t.Errorf("Expected state %s, got %s", PlayerStateOffline, player.State)
// 	}
// 	if player.MatchHistory == nil {
// 		t.Error("Expected MatchHistory to be initialized")
// 	}
// }

// func TestPlayerState_IsValid(t *testing.T) {
// 	tests := []struct {
// 		state PlayerState
// 		valid bool
// 	}{
// 		{PlayerStateInGame, true},
// 		{PlayerStateInQuery, true},
// 		{PlayerStateOffline, true},
// 		{PlayerStateOnline, true},
// 		{PlayerState("Invalid"), false},
// 	}
	
// 	for _, test := range tests {
// 		if test.state.IsValid() != test.valid {
// 			t.Errorf("Expected %s to be valid: %v", test.state, test.valid)
// 		}
// 	}
// }

// func TestPlayer_UpdateLevel(t *testing.T) {
// 	player := NewPlayer("test", "pass", "email@test.com")
// 	oldTime := player.UpdatedAt
	
// 	time.Sleep(time.Millisecond) // Ensure time difference
// 	player.UpdateLevel(5.5)
	
// 	if player.Level != 5.5 {
// 		t.Errorf("Expected level 5.5, got %f", player.Level)
// 	}
// 	if !player.UpdatedAt.After(oldTime) {
// 		t.Error("Expected UpdatedAt to be updated")
// 	}
// }

// func TestPlayer_AddMatch(t *testing.T) {
// 	player := NewPlayer("test", "pass", "email@test.com")
	
// 	player.AddMatch("match1", true)
// 	player.AddMatch("match2", false)
	
// 	if len(player.MatchHistory) != 2 {
// 		t.Errorf("Expected 2 matches, got %d", len(player.MatchHistory))
// 	}
// 	if !player.MatchHistory["match1"] {
// 		t.Error("Expected match1 to be won")
// 	}
// 	if player.MatchHistory["match2"] {
// 		t.Error("Expected match2 to be lost")
// 	}
// }

// func TestPlayer_SetState(t *testing.T) {
// 	player := NewPlayer("test", "pass", "email@test.com")
	
// 	err := player.SetState(PlayerStateOnline)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}
// 	if player.State != PlayerStateOnline {
// 		t.Errorf("Expected state %s, got %s", PlayerStateOnline, player.State)
// 	}
	
// 	err = player.SetState(PlayerState("Invalid"))
// 	if err == nil {
// 		t.Error("Expected error for invalid state")
// 	}
// }

// func TestPlayer_GetWinRate(t *testing.T) {
// 	player := NewPlayer("test", "pass", "email@test.com")
	
// 	// No matches
// 	if player.GetWinRate() != 0.0 {
// 		t.Errorf("Expected win rate 0.0, got %f", player.GetWinRate())
// 	}
	
// 	// Add matches
// 	player.AddMatch("match1", true)
// 	player.AddMatch("match2", true)
// 	player.AddMatch("match3", false)
	
// 	expected := 2.0 / 3.0
// 	if player.GetWinRate() != expected {
// 		t.Errorf("Expected win rate %f, got %f", expected, player.GetWinRate())
// 	}
// }

// func TestPlayer_IsOnline(t *testing.T) {
// 	player := NewPlayer("test", "pass", "email@test.com")
	
// 	// Offline
// 	if player.IsOnline() {
// 		t.Error("Expected player to be offline")
// 	}
	
// 	// Online states
// 	onlineStates := []PlayerState{PlayerStateOnline, PlayerStateInGame, PlayerStateInQuery}
// 	for _, state := range onlineStates {
// 		player.SetState(state)
// 		if !player.IsOnline() {
// 			t.Errorf("Expected player to be online when state is %s", state)
// 		}
// 	}
// }

// func TestPlayer_Validate(t *testing.T) {
// 	// Valid player
// 	player := NewPlayer("testuser", "password", "test@example.com")
// 	if err := player.Validate(); err != nil {
// 		t.Errorf("Unexpected validation error: %v", err)
// 	}
	
// 	// Invalid name
// 	player.Name = ""
// 	if err := player.Validate(); err == nil {
// 		t.Error("Expected validation error for empty name")
// 	}
	
// 	player.Name = "ab" // Too short
// 	if err := player.Validate(); err == nil {
// 		t.Error("Expected validation error for short name")
// 	}
	
// 	// Invalid level
// 	player.Name = "validname"
// 	player.Level = -1
// 	if err := player.Validate(); err == nil {
// 		t.Error("Expected validation error for negative level")
// 	}
// }

// func TestPlayer_JSON(t *testing.T) {
// 	player := NewPlayer("test", "password", "test@example.com")
// 	player.ID = "123"
	
// 	// Marshal
// 	data, err := json.Marshal(player)
// 	if err != nil {
// 		t.Fatalf("Failed to marshal player: %v", err)
// 	}
	
// 	// Check password is not included
// 	var result map[string]interface{}
// 	json.Unmarshal(data, &result)
// 	if _, exists := result["password"]; exists {
// 		t.Error("Password should not be included in JSON")
// 	}
	
// 	// Unmarshal
// 	var newPlayer Player
// 	err = json.Unmarshal(data, &newPlayer)
// 	if err != nil {
// 		t.Fatalf("Failed to unmarshal player: %v", err)
// 	}
	
// 	if newPlayer.ID != player.ID {
// 		t.Errorf("Expected ID %s, got %s", player.ID, newPlayer.ID)
// 	}
// 	if newPlayer.Name != player.Name {
// 		t.Errorf("Expected name %s, got %s", player.Name, newPlayer.Name)
// 	}
// }