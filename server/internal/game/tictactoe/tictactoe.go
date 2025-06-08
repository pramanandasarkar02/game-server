package tictactoe

import (
	"context"
	"fmt"
	"sync"

	"github.com/pramanandasarkar02/game-server/pkg/logger"
)

type TicTacToe struct {
	states map[string]*TicTacToeState
	mutex  sync.RWMutex
}

type Move struct {
	Index    int    `json:"index"`
	PlayerID string `json:"playerID"`
}

func NewTicTacToe() *TicTacToe {
	return &TicTacToe{
		states: make(map[string]*TicTacToeState),
	}
}

func (t *TicTacToe) ID() string {
	return "a2"
}

func (t *TicTacToe) Title() string {
	return "Tic Tac Toe"
}

func (t *TicTacToe) RequiredPlayers() int {
	return 2
}

func (t *TicTacToe) InitializeState(matchID string, players []string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.states[matchID] = &TicTacToeState{
		Board:   [9]string{},
		Turn:    players[0],
		Players: players,
	}
}

func (t *TicTacToe) GetState(matchID string) interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	if state, exists := t.states[matchID]; exists {
		return state
	}
	return nil
}

func (t *TicTacToe) HandleMove(ctx context.Context, matchID string, playerID string, moveData interface{}) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	state, exists := t.states[matchID]
	if !exists {
		return fmt.Errorf("match %s not found", matchID)
	}
	if state.Winner != "" || state.IsDraw {
		return fmt.Errorf("game is over")
	}
	if state.Turn != playerID {
		return fmt.Errorf("not player's turn")
	}

	move, ok := moveData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid move format")
	}
	indexFloat, ok := move["index"].(float64)
	if !ok || indexFloat < 0 || indexFloat > 8 {
		return fmt.Errorf("invalid move index")
	}
	index := int(indexFloat)
	if state.Board[index] != "" {
		return fmt.Errorf("cell already occupied")
	}

	symbol := "X"
	if state.Players[1] == playerID {
		symbol = "O"
	}
	state.Board[index] = symbol

	if winner, gameOver := checkWin(state.Board); gameOver {
		if winner != "" {
			state.Winner = playerID
		} else {
			state.IsDraw = true
		}
	} else {
		state.Turn = state.Players[0]
		if state.Turn == playerID {
			state.Turn = state.Players[1]
		}
	}

	logger.Info("Move processed for match %s by player %s: index %d", matchID, playerID, index)
	return nil
}

func checkWin(board [9]string) (string, bool) {
	winPatterns := [][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // Rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // Columns
		{0, 4, 8}, {2, 4, 6}, // Diagonals
	}
	for _, pattern := range winPatterns {
		a, b, c := pattern[0], pattern[1], pattern[2]
		if board[a] != "" && board[a] == board[b] && board[a] == board[c] {
			return board[a], true
		}
	}
	if !contains(board[:], "") {
		return "", true
	}
	return "", false
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
