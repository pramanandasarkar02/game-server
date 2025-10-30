package snake

import (
	"math/rand/v2"
	"sync"
)

type Food struct {
	Position Point `json:"position"`
	Value    int   `json:"value"`
}

type Obstacle struct {
	Object []Point `json:"object"`
}

func CreateNewRandomObstacle(width, height int) Obstacle {
	startingPoint := Point{
		X: rand.IntN(width),
		Y: rand.IntN(height),
	}

	length := 3 + rand.IntN(4)

	object := make([]Point, 0, length)
	object = append(object, startingPoint)

	for i := 1; i < length; {
		dir := rand.IntN(4)
		newPoint := startingPoint
		switch dir {
		case 0:
			newPoint.X += 1
		case 1:
			newPoint.X -= 1
		case 2:
			newPoint.Y += 1
		case 3:
			newPoint.Y -= 1
		}

		if newPoint.X >= 0 && newPoint.X < width && newPoint.Y >= 0 && newPoint.Y < height {
			i += 1
			object = append(object, newPoint)
			startingPoint = newPoint
		}
	}
	return Obstacle{
		Object: object,
	}
}

type SnakeBoard struct {
	SnakeControllers map[string]*SnakeController
	Foods            []Food
	Width            int
	Height           int
	minimumFood      int
	numberOfFoodRange int
	obstacleCount    int
	Obstacles        []Obstacle
	mu               sync.RWMutex
}

type SnakeBoardPlayerInformation struct {
	PlayerId    string     `json:"playerId"`
	PlayerSnake Snake      `json:"playerSnake"`
	Foods       []Food     `json:"foods"`
	OtherSnakes []Snake    `json:"otherSnakes"`
	Obstacles   []Obstacle `json:"obstacles"`
}

func NewSnakeBoard() *SnakeBoard {
	snakeControllers := make(map[string]*SnakeController)
	height := 40
	width := 60

	obsCount, obstacles := createObstacles(width, height)
	snakeBoard := &SnakeBoard{
		SnakeControllers:  snakeControllers,
		Foods:             make([]Food, 0),
		Obstacles:         obstacles,
		Width:             width,
		Height:            height,
		minimumFood:       4,
		numberOfFoodRange: 3,
		obstacleCount:     obsCount,
	}
	snakeBoard.GenerateFood()
	return snakeBoard
}

func (sb *SnakeBoard) AddPlayer(playerId string) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	if _, exists := sb.SnakeControllers[playerId]; !exists {
		sb.SnakeControllers[playerId] = NewSnakeController(NewSnake())
	}
}

func (sb *SnakeBoard) GenerateFood() {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	snakes := make([]Snake, 0)
	for _, sc := range sb.SnakeControllers {
		snakes = append(snakes, *sc.Snake)
	}
	
	numberOfFood := sb.minimumFood + rand.IntN(sb.numberOfFoodRange)
	for len(sb.Foods) < numberOfFood {
		x := rand.IntN(sb.Width)
		y := rand.IntN(sb.Height)

		newFood := Food{
			Position: Point{
				X: x,
				Y: y,
			},
			Value: 1 + rand.IntN(5),
		}
		if sb.isOccupied(newFood.Position, snakes, sb.Obstacles) {
			continue
		}
		sb.Foods = append(sb.Foods, newFood)
	}
}

func (sb *SnakeBoard) isOccupied(p Point, snakes []Snake, obstacles []Obstacle) bool {
	for _, s := range snakes {
		for _, body := range s.SnakeBody {
			if body.X == p.X && body.Y == p.Y {
				return true
			}
		}
		if s.SnakeHead.X == p.X && s.SnakeHead.Y == p.Y {
			return true
		}
	}

	for _, obs := range obstacles {
		for _, o := range obs.Object {
			if o.X == p.X && o.Y == p.Y {
				return true
			}
		}
	}

	return false
}

func (sb *SnakeBoard) ExecutePlayerMovement(playerId string, direction Direction) {
	sb.mu.RLock()
	sc, ok := sb.SnakeControllers[playerId]
	sb.mu.RUnlock()

	if ok {
		sc.KeyboardController(direction)
	}
}

func (sb *SnakeBoard) RunSnake(playerId string) (bool, string) {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	sc, ok := sb.SnakeControllers[playerId]
	if !ok {
		return false, "player not found"
	}

	isCol, msg := sc.RunSnake(sb)
	return isCol, msg
}

func (sb *SnakeBoard) GetSnakeBoard(playerId string) *SnakeBoardPlayerInformation {
	sb.mu.RLock()
	defer sb.mu.RUnlock()

	snakeController, ok := sb.SnakeControllers[playerId]
	if !ok {
		return &SnakeBoardPlayerInformation{
			PlayerId:    playerId,
			Foods:       sb.Foods,
			Obstacles:   sb.Obstacles,
			OtherSnakes: make([]Snake, 0),
		}
	}

	playerSnake := snakeController.Snake
	foods := sb.Foods
	obstacles := sb.Obstacles
	otherSnakes := make([]Snake, 0)
	
	for pId, sc := range sb.SnakeControllers {
		if pId != playerId {
			otherSnakes = append(otherSnakes, *sc.Snake)
		}
	}

	return &SnakeBoardPlayerInformation{
		PlayerId:    playerId,
		PlayerSnake: *playerSnake,
		Foods:       foods,
		Obstacles:   obstacles,
		OtherSnakes: otherSnakes,
	}
}

func createObstacles(w, h int) (int, []Obstacle) {
	minimumObstacles := 2
	numberOfObstacleRange := 3

	obstacleCount := minimumObstacles + rand.IntN(numberOfObstacleRange)

	obstacles := make([]Obstacle, 0, obstacleCount)
	for i := 0; i < obstacleCount; i++ {
		obstacles = append(obstacles, CreateNewRandomObstacle(w, h))
	}

	return obstacleCount, obstacles
}