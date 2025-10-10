package snake


const(
	MAX_RUNNING_SNAKE_GAMES = 10
)

type SnakeService struct{
	SnakeBoard []SnakeBoard
}



func NewSnakeService() *SnakeService{
	// comsume full space for the game thats bad 
	return &SnakeService{
		SnakeBoard: make([]SnakeBoard, MAX_RUNNING_SNAKE_GAMES),
	}
}