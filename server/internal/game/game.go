package game

import (
	"context"
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

var games = make(map[string]Game)

func RegisterGame(g Game) {
	games[g.ID()] = g
}

func GetGame(id string) Game {
	return games[id]
}

func TicTacToeGame() Game {
	return tictactoe.NewTicTacToe()
}