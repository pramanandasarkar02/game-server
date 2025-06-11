package game

import (
    "context"
    "sync"

    "github.com/pramanandasarkar02/game-server/internal/game/snake"
    "github.com/pramanandasarkar02/game-server/internal/game/tictactoe"
)

type Game interface {
    ID() string
    Title() string
    RequiredPlayers() int
    HandleMove(ctx context.Context, matchID string, playerID string, move interface{}) error
    GetState(matchID string) interface{}
    InitializeState(matchID string, players []string)
}

type gameRegistry struct {
    games map[string]Game
    mutex sync.RWMutex
}

var registry = &gameRegistry{
    games: make(map[string]Game),
}

func init() {
    RegisterGame(TicTacToeGame())
    RegisterGame(SnakeGame())
}

func RegisterGame(g Game) {
    registry.mutex.Lock()
    defer registry.mutex.Unlock()
    registry.games[g.ID()] = g
}

func GetGame(id string) Game {
    registry.mutex.RLock()
    defer registry.mutex.RUnlock()
    return registry.games[id]
}

func TicTacToeGame() Game {
    return tictactoe.NewTicTacToe()
}

func SnakeGame() Game {
    return snake.NewSnakeGame()
}