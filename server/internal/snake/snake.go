package snake

import (
	"log"
	"time"
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Score struct {
	Value int `json:"value"`
}

type Direction string

const (
	LEFT  Direction = "LEFT"
	RIGHT Direction = "RIGHT"
	UP    Direction = "UP"
	DOWN  Direction = "DOWN"
)

type Snake struct {
	SnakeHead    Point     `json:"snakeHead"`
	SnakeBody    []Point   `json:"snakeBody"`
	Direction    Direction `json:"direction"`
	Score        Score     `json:"score"`
	StartingTime time.Time `json:"time"`
}

func NewSnake() *Snake {
	return &Snake{
		SnakeHead:    Point{X: 5, Y: 5},
		SnakeBody:    []Point{{X: 4, Y: 5}, {X: 3, Y: 5}},
		Direction:    RIGHT,
		Score:        Score{Value: 0},
		StartingTime: time.Now(),
	}
}

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

func executeDirMovement(head Point, dir Direction) Point {
	switch dir {
	case LEFT:
		head.X -= 1
	case RIGHT:
		head.X += 1
	case UP:
		head.Y -= 1
	case DOWN:
		head.Y += 1
	}
	return head
}

func checkFood(foods []Food, head Point) (bool, Food, int) {
	for idx, food := range foods {
		if food.Position.X == head.X && food.Position.Y == head.Y {
			return true, food, idx
		}
	}
	return false, Food{}, -1
}

func checkCollision(head Point, snakeBody []Point, gameBoard *SnakeBoard) (bool, string) {
	// Board range check
	if (head.X < 0 || head.X >= gameBoard.Width) || (head.Y < 0 || head.Y >= gameBoard.Height) {
		return true, "Out of range"
	}

	// Obstacle collision
	for _, obs := range gameBoard.Obstacles {
		for _, o := range obs.Object {
			if head.X == o.X && head.Y == o.Y {
				return true, "Hit Obstacle"
			}
		}
	}

	// Other snakes collision
	for _, sc := range gameBoard.SnakeControllers {
		otherSnake := sc.Snake
		if otherSnake.SnakeHead.X == head.X && otherSnake.SnakeHead.Y == head.Y {
			return true, "Hit other snake head"
		}
		for _, part := range otherSnake.SnakeBody {
			if head.X == part.X && head.Y == part.Y {
				return true, "Hit other snake body"
			}
		}
	}

	// Self body collision
	for _, part := range snakeBody {
		if head.X == part.X && head.Y == part.Y {
			return true, "Self Collision"
		}
	}

	return false, ""
}

func executeMovement(newHead Point, snake *Snake, isFood bool) {
	snake.SnakeBody = append(snake.SnakeBody, snake.SnakeHead)
	if !isFood && len(snake.SnakeBody) > 0 {
		snake.SnakeBody = snake.SnakeBody[1:]
	}
	snake.SnakeHead = newHead
}

func (s *Snake) Movement(gameBoard *SnakeBoard) (bool, string) {
	newHeadPosition := executeDirMovement(s.SnakeHead, s.Direction)

	if isCollision, msg := checkCollision(newHeadPosition, s.SnakeBody, gameBoard); isCollision {
		log.Printf("Collision detected: %s", msg)
		return true, msg
	}

	if isFood, food, idx := checkFood(gameBoard.Foods, newHeadPosition); isFood {
		executeMovement(newHeadPosition, s, true)
		s.Score.Value += food.Value
		
		// Remove eaten food
		if idx >= 0 && idx < len(gameBoard.Foods) {
			gameBoard.Foods = append(gameBoard.Foods[:idx], gameBoard.Foods[idx+1:]...)
		}
		return false, ""
	}

	executeMovement(newHeadPosition, s, false)
	return false, ""
}