package snake

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"
)

type SnakeGameBoard struct {
    Board   [8][8]string `json:"board"`
    Turn    string       `json:"turn"`
    Winner  string       `json:"winner"`
    Players []string     `json:"players"`
}

type SnakeGame struct {
    boards map[string]*SnakeGameBoard
    mutex  sync.RWMutex
}

type Move struct {
    X int `json:"x"`
    Y int `json:"y"`
}

func NewSnakeGame() *SnakeGame {
    return &SnakeGame{
        boards: make(map[string]*SnakeGameBoard),
    }
}

func (s *SnakeGame) ID() string {
    return "a4"
}

func (s *SnakeGame) Title() string {
    return "Snake Game"
}

func (s *SnakeGame) RequiredPlayers() int {
    return 4
}

func (s *SnakeGame) InitializeState(matchID string, players []string) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    board := &SnakeGameBoard{
        Board:   [8][8]string{},
        Turn:    players[0],
        Players: players,
    }
    board.Board[0][0] = players[0]
    board.Board[7][0] = players[1]
    board.Board[7][7] = players[2]
    board.Board[0][7] = players[3]
    s.boards[matchID] = board
}

func (s *SnakeGame) GetState(matchID string) interface{} {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    board, exists := s.boards[matchID]
    if !exists {
        return nil
    }
    return board
}

func checkValidMove(board [8][8]string, x, y int, player string) bool {
    if x < 0 || x >= 8 || y < 0 || y >= 8 {
        return false
    }
    if board[x][y] != "" {
        return false
    }
    for dx := -1; dx <= 1; dx++ {
        for dy := -1; dy <= 1; dy++ {
            if dx == 0 && dy == 0 {
                continue
            }
            nx, ny := x + dx, y + dy
            if nx >= 0 && nx < 8 && ny >= 0 && ny < 8 && board[nx][ny] == player {
                return true
            }
        }
    }
    return false
}

func checkTermination(board [8][8]string, player string) bool {
    for x := 0; x < 8; x++ {
        for y := 0; y < 8; y++ {
            if checkValidMove(board, x, y, player) {
                return false
            }
        }
    }
    return true
}

func (s *SnakeGame) HandleMove(ctx context.Context, matchID string, playerID string, moveData interface{}) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    board, exists := s.boards[matchID]
    if !exists {
        return fmt.Errorf("match %s not found", matchID)
    }

    if board.Winner != "" {
        return fmt.Errorf("game is over")
    }
    if board.Turn != playerID {
        return fmt.Errorf("not player's turn")
    }

    var move Move
    data, err := json.Marshal(moveData)
    if err != nil {
        return fmt.Errorf("failed to marshal move data: %v", err)
    }
    if err := json.Unmarshal(data, &move); err != nil {
        return fmt.Errorf("invalid move data: %v", err)
    }

    if !checkValidMove(board.Board, move.X, move.Y, playerID) {
        return fmt.Errorf("invalid move to (%d, %d)", move.X, move.Y)
    }

    board.Board[move.X][move.Y] = playerID

    if checkTermination(board.Board, playerID) {
        board.Winner = playerID
        return nil
    }

    currentIndex := 0
    for i, p := range board.Players {
        if p == playerID {
            currentIndex = i
            break
        }
    }
    nextIndex := (currentIndex + 1) % len(board.Players)
    board.Turn = board.Players[nextIndex]

    return nil
}