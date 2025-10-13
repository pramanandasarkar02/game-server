package snake

import "game-server/internal/service"


const(
	MAX_RUNNING_SNAKE_GAMES = 10
	BOARD_WIDTH = 60
	BOARD_HEIGHT = 40
	CELL_SIZE = 10
)

type SnakeService struct{
	SnakeBoard []SnakeBoard
}


type SnakeGameMetaDataResponse struct {
	BoardWidth  int `json:"boardWidth"`
	BoardHeight int `json:"boardHeight"`
	CellSize    int `json:"cellSize"`
}



func NewSnakeService() *SnakeService{
	// comsume full space for the game thats bad 
	return &SnakeService{
		SnakeBoard: make([]SnakeBoard, MAX_RUNNING_SNAKE_GAMES),
	}
}



func(ss * SnakeService) SnakeGameMetaData() *SnakeGameMetaDataResponse {
	
	return &SnakeGameMetaDataResponse{
		BoardWidth:  BOARD_WIDTH,
		BoardHeight: BOARD_HEIGHT,
		CellSize:    CELL_SIZE,
	}
}

func (ss *SnakeService) StartGame(gameEnv service.GameEnv ){
	
}