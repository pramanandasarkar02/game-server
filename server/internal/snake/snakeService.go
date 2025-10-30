package snake

import (
	"log"
	"sync"
)

const (
	MAX_RUNNING_SNAKE_GAMES = 10
	BOARD_WIDTH             = 60
	BOARD_HEIGHT            = 40
	CELL_SIZE               = 10
)

type SnakeService struct {
	SnakeBoards  map[string]*SnakeBoard
	MatchPlayers map[string][]string
	mu           sync.RWMutex
}

type SnakeGameMetaDataResponse struct {
	BoardWidth  int `json:"boardWidth"`
	BoardHeight int `json:"boardHeight"`
	CellSize    int `json:"cellSize"`
}

func NewSnakeService() *SnakeService {
	return &SnakeService{
		SnakeBoards:  make(map[string]*SnakeBoard),
		MatchPlayers: make(map[string][]string),
	}
}

func (ss *SnakeService) SnakeGameMetaData() *SnakeGameMetaDataResponse {
	return &SnakeGameMetaDataResponse{
		BoardWidth:  BOARD_WIDTH,
		BoardHeight: BOARD_HEIGHT,
		CellSize:    CELL_SIZE,
	}
}

func (ss *SnakeService) StartGame(matchId string, playerIds []string) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if _, ok := ss.SnakeBoards[matchId]; ok {
		return
	}
	ss.SnakeBoards[matchId] = NewSnakeBoard()
	ss.MatchPlayers[matchId] = playerIds
}

func (ss *SnakeService) AddPlayer(matchId, playerId string) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if sb, ok := ss.SnakeBoards[matchId]; ok {
		sb.AddPlayer(playerId)
	}
}

func (ss *SnakeService) ExecuteMovement(matchId, playerId string, direction Direction) {
	ss.mu.RLock()
	snakeBoard, ok := ss.SnakeBoards[matchId]
	ss.mu.RUnlock()

	if ok {
		snakeBoard.ExecutePlayerMovement(playerId, direction)
	}
}

func (ss *SnakeService) GenerateFood(matchId string) {
	ss.mu.RLock()
	sb, ok := ss.SnakeBoards[matchId]
	ss.mu.RUnlock()

	if ok {
		sb.GenerateFood()
	}
}

func (ss *SnakeService) GetBoardStats(matchId, playerId string) *SnakeBoardPlayerInformation {
	ss.mu.RLock()
	sb, ok := ss.SnakeBoards[matchId]
	ss.mu.RUnlock()

	if ok {
		return sb.GetSnakeBoard(playerId)
	}
	return &SnakeBoardPlayerInformation{}
}

func (ss *SnakeService) RunAllSnake(matchId string) {
	ss.mu.RLock()
	sb, ok := ss.SnakeBoards[matchId]
	players := ss.MatchPlayers[matchId]
	ss.mu.RUnlock()

	if !ok {
		log.Println("Match id not found")
		return
	}

	for _, playerId := range players {
		sb.RunSnake(playerId)
	}
}

func (ss *SnakeService) RunSnake(matchId, playerId string) (bool, string) {
	ss.mu.RLock()
	sb, ok := ss.SnakeBoards[matchId]
	ss.mu.RUnlock()

	if !ok {
		return false, "snake board not found"
	}
	return sb.RunSnake(playerId)
}

func (ss *SnakeService) EndGame(matchId string) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	delete(ss.SnakeBoards, matchId)
	delete(ss.MatchPlayers, matchId)
}