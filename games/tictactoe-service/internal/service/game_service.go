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