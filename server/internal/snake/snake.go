package snake

import (
	"log"
	"time"
)

// Point represents a position on the grid
type Point struct {
	X int	`json:"x"`
	Y int	`json:"y"`
}

// Score holds the current score of the snake
type Score struct {
	Value int 	`json:"value"`
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
	SnakeHead    Point 	`json:"snakeHead"`
	SnakeBody    []Point `json:"snakeBody"`
	Direction    Direction	`json:"direction"`
	Score        Score		`json:"score"`
	StartingTime time.Time	`json:"time"`
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
func (s *Snake) Controller(c Direction) SnakeControllerResponse {
	switch c {
	case UP:
		if s.Direction == LEFT || s.Direction == RIGHT {
			s.Direction = UP
			return SnakeControllerResponse{true, "Direction changed to UP"}
		}
	case DOWN:
		if s.Direction == LEFT || s.Direction == RIGHT {
			s.Direction = DOWN
			return SnakeControllerResponse{true, "Direction changed to DOWN"}
		}
	case LEFT:
		if s.Direction == UP || s.Direction == DOWN {
			s.Direction = LEFT
			return SnakeControllerResponse{true, "Direction changed to LEFT"}
		}
	case RIGHT:
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
		if food.Position.X == head.X && food.Position.Y == head.Y {
			return true, food
		}
	}
	return false, Food{}
}

func checkCollision(head Point, snakeBody []Point, gameBoard *SnakeBoard)(bool, string) {
	// board reange
	if (head.X < 0 || head.X >= gameBoard.Width) || (head.Y < 0 || head.Y >= gameBoard.Height) {
		return true, "Out of range"
	}

	// obstacle
	// for _, obs := range gameBoard.obstacles {
	// 	if head.X == obs.X && head.Y == obs.Y{
	// 		return true, "Hit Obstacle"
	// 	}
	// }

	// other snake 
	

	// self body
	for _, part := range snakeBody {
		if head.X == part.X && head.X == part.Y {
			return true, "Self Collision"
		}
	}
	
	return false, ""
}

func executeMovement(newHead Point, snake *Snake, isFood bool){
	snake.SnakeBody = append(snake.SnakeBody, snake.SnakeHead)
	if len(snake.SnakeBody) > 0 && !isFood{
		snake.SnakeBody = snake.SnakeBody[1:]
	}
	snake.SnakeHead = newHead
}

// One tick time one execution
func (s *Snake) Movement(gameBoard *SnakeBoard) {
	newHeadPosition := exeucteDirMovement(s.SnakeHead, s.Direction)

	if isCollision, msg := checkCollision(newHeadPosition, s.SnakeBody, gameBoard); isCollision {
		log.Panicf("there is a collision %s", msg)
	}

	if isFood, food := checkFood(gameBoard.Foods, newHeadPosition); isFood {
		executeMovement(newHeadPosition, s, isFood)
		s.Score.Value += food.Value
	}

	executeMovement(newHeadPosition, s, false)
}
