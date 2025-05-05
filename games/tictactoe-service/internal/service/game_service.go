package service

import (
	"fmt"
	"sync"

	"github.com/pramanandasarkar02/game-server/games/tictactoe-service/internal/model"
)





type GameService struct {
	games map[string]*model.Game
	mutex sync.Mutex	
}


func NewGameService() *GameService {
	return &GameService{
		games: make(map[string]*model.Game),
	}
}



func (g *GameService) CreateGame() *model.Game {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	id := fmt.Sprintf("tic_tac_toe-%d", len(g.games) + 1)
	game := &model.Game{
		Id: id,
		Board: [3][3]string{},
		CurrentPlayer: "X",
		Status: "IN_PROGRESS",
		Winner: "",
	}

	g.games[id] = game
	return game
}
func (g * GameService) GetGame(id string) (*model.Game, bool) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	game, exists := g.games[id]
	return game, exists
}

func (g *GameService) MakeMove(id string, x int, y int, player string) (*model.Game, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	game, exists := g.games[id]	
	if !exists {
		return nil, fmt.Errorf("game not found")
	}

	if game.Status != "IN_PROGRESS" {
		return nil, fmt.Errorf("game is over")
	}

	if game.Board[x][y] != "" {
		return nil, fmt.Errorf("invalid move")
	}	

	game.Board[x][y] = player

	if game.Board[0][0] == player && game.Board[0][1] == player && game.Board[0][2] == player {
		game.Status = "WINNER"
		game.Winner = player
	} else if game.Board[1][0] == player && game.Board[1][1] == player && game.Board[1][2] == player {
		game.Status = "WINNER"
		game.Winner = player
	} else if game.Board[2][0] == player && game.Board[2][1] == player && game.Board[2][2] == player {
		game.Status = "WINNER"
		game.Winner = player
	} else if game.Board[0][0] == player && game.Board[1][0] == player && game.Board[2][0] == player {
		game.Status = "WINNER"
		game.Winner = player
	} else if game.Board[0][1] == player && game.Board[1][1] == player && game.Board[2][1] == player {
		game.Status = "WINNER"
		game.Winner = player	
	} else if game.Board[0][2] == player && game.Board[1][2] == player && game.Board[2][2] == player {
		game.Status = "WINNER"
		game.Winner = player
	} else if game.Board[0][0] == player && game.Board[1][1] == player && game.Board[2][2] == player {
		game.Status = "WINNER"
		game.Winner = player
	} else if game.Board[0][2] == player && game.Board[1][1] == player && game.Board[2][0] == player {
		game.Status = "WINNER"
		game.Winner = player
	}

	if game.Status == "WINNER" {
		game.CurrentPlayer = "X"
	} else {
		if game.CurrentPlayer == "X" {
			game.CurrentPlayer = "O"
		} else {
			game.CurrentPlayer = "X"
		}
	}
	return game, nil
}