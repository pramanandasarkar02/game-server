package snake

import "time"

// Point represents a position on the grid
type Point struct {
	X int32
	Y int32
}

// Score holds the current score of the snake
type Score struct {
	Value int
}

// Direction type for snake movement
type Direction string

const (
	LEFT  Direction = "left"
	RIGHT Direction = "right"
	UP    Direction = "up"
	DOWN  Direction = "down"
)

// Snake represents the snake entity
type Snake struct {
	SnakeHead    Point
	SnakeBody    []Point
	Direction    Direction
	Score        Score
	StartingTime time.Time
}

// NewSnake initializes and returns a new Snake instance
func NewSnake() *Snake {
	return &Snake{
		SnakeHead:    Point{X: 5, Y: 5},
		SnakeBody:    []Point{{X: 4, Y: 5}, {X: 3, Y: 5}},
		Direction:    RIGHT,
		Score:        Score{Value: 0},
		StartingTime: time.Now(),
	}
}

// snake direction controller
func (s *Snake) Controller(c ControllerOption) SnakeControllerResponse {
	switch c {
	case UP_CONTROLLER:
		if s.Direction == LEFT || s.Direction == RIGHT {
			s.Direction = UP
			return SnakeControllerResponse{true, "Direction changed to UP"}
		}
	case DOWN_CONTROLLER:
		if s.Direction == LEFT || s.Direction == RIGHT {
			s.Direction = DOWN
			return SnakeControllerResponse{true, "Direction changed to DOWN"}
		}
	case LEFT_CONTROLLER:
		if s.Direction == UP || s.Direction == DOWN {
			s.Direction = LEFT
			return SnakeControllerResponse{true, "Direction changed to LEFT"}
		}
	case RIGHT_CONTROLLER:
		if s.Direction == UP || s.Direction == DOWN {
			s.Direction = RIGHT
			return SnakeControllerResponse{true, "Direction changed to RIGHT"}
		}
	}

	return SnakeControllerResponse{false, "Invalid direction change"}
}

func exeucteDirMovement(head Point, dir Direction) Point {
	switch dir {
	case LEFT:
		head.X += 1
	case RIGHT:
		head.X -= 1
	case UP:
		head.Y -= 1
	case DOWN:
		head.Y += 1
	}
	return head

}

func checkFood(foods []Food, head Point) (bool, Food) {
	for _, food := range foods {
		if food.Position.X == head.X && food.Position.X == head.Y {
			return true, food
		}
	}
	return false, Food{}
}

func checkCollision(head Point, snakes []Snake) {

}

func executeMovement(newHead Point, snake *Snake, isFood bool){
	snake.SnakeBody = append(snake.SnakeBody, snake.SnakeHead)
	if len(snake.SnakeBody) > 0 && !isFood{
		snake.SnakeBody = snake.SnakeBody[1:]
	}
	snake.SnakeHead = newHead
}

// One tick time one execution
func (s *Snake) Movement(gameBoard SnakeBoard) {
	newHeadPosition := exeucteDirMovement(s.SnakeHead, s.Direction)
	if isFood, food := checkFood(gameBoard.Foods, newHeadPosition); isFood {
		executeMovement(newHeadPosition, s, isFood)
		s.Score.Value += food.Value
	}

}
