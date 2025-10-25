package snake



const(
	MAX_RUNNING_SNAKE_GAMES = 10
	BOARD_WIDTH = 60
	BOARD_HEIGHT = 40
	CELL_SIZE = 10
)

type SnakeService struct{
	SnakeBoards map[string]*SnakeBoard // matchId -> snakeboard
	MatchPlayers map[string][]string // matchId -> playerId
}


type SnakeGameMetaDataResponse struct {
	BoardWidth  int `json:"boardWidth"`
	BoardHeight int `json:"boardHeight"`
	CellSize    int `json:"cellSize"`
}



func NewSnakeService() *SnakeService{
	// comsume full space for the game thats bad 
	return &SnakeService{
		SnakeBoards :make(map[string]*SnakeBoard),
		MatchPlayers: make(map[string][]string),
	}
}



func(ss * SnakeService) SnakeGameMetaData() *SnakeGameMetaDataResponse {
	return &SnakeGameMetaDataResponse{
		BoardWidth:  BOARD_WIDTH,
		BoardHeight: BOARD_HEIGHT,
		CellSize:    CELL_SIZE,
	}
}

func (ss *SnakeService) StartGame(matchId string, playerIds []string ){
	
	if _, ok := ss.SnakeBoards[matchId]; ok{
		return 
	}
	ss.SnakeBoards[matchId] = NewSnakeBoard()
	ss.MatchPlayers[matchId] = playerIds
	
}

func(ss* SnakeService) AddPlayer(matchId, playerId string){
	ss.SnakeBoards[matchId].AddPlayer(playerId)
}

func (ss *SnakeService) ExecuteMovement(matchId, playerId string, direction Direction){
	// ss.SnakeBoard
	snakeBoard := ss.SnakeBoards[matchId]
	snakeBoard.ExecutePlayerMovement(playerId, direction)
}

func(ss *SnakeService) GenerateFood(matchId string){
	if sb, ok := ss.SnakeBoards[matchId]; ok{
		sb.GenerateFood()
	}
}

func(ss *SnakeService)GetBoardStats(matchId, playerId string) *SnakeBoardPlayerInformation{
	if sb, ok := ss.SnakeBoards[matchId]; ok{
		return sb.GetSnakeBoard(playerId)
	}
	return &SnakeBoardPlayerInformation{}
}
func (ss *SnakeService)RunAllSnake(matchId string){
	
}

func(ss *SnakeService)RunSnake(matchId, playerId string)(bool, string){
	if sb, ok := ss.SnakeBoards[matchId]; ok {
		return sb.RunSnake(playerId)
	}
	return true, ""
}

func (ss *SnakeService) EndGame(){
	
}