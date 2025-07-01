package models

import (
    "testing"
)

func TestNewGame(t *testing.T) {
    title := "Tic-Tac-Toe"
    requiredPlayer := 2
    metaData := "{'requiredPlayer': 2}"
    game := NewGame(title, requiredPlayer, metaData)

    if game == nil {
        t.Fatal("Expected non-nil game")
    }
    if game.ID == "" {
        t.Fatal("Expected non-empty game ID")
    }
}