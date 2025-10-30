package snake

import "sync"

type SnakeController struct {
	Snake *Snake
	mu    sync.Mutex
}

type SnakeControllerResponse struct {
	Ok  bool
	Msg string
}

func NewSnakeController(snake *Snake) *SnakeController {
	return &SnakeController{
		Snake: snake,
	}
}

func (sc *SnakeController) RunSnake(snakeBoard *SnakeBoard) (bool, string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	isCol, msg := sc.Snake.Movement(snakeBoard)
	return isCol, msg
}

func (sc *SnakeController) KeyboardController(option Direction) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.Snake.Controller(option)
	return nil
}