package snake

import (
	"math/rand/v2"
)

type Food struct{
	Position Point 		`json:"position"`
	Value int			`json:"value"`
}

type Obstacle struct{
	Object []Point 		`json:"object"`
}


type SnakeBoard struct{
	SnakeControllers map[string]*SnakeController  // playerId -> snakecontroller
	Foods []Food
	Width int 
	Height int
	minimumFood int 
	numberOfFoodRange int 
	Obstacles []Obstacle
}

type SnakeBoardPlayerInformation struct{
	PlayerId string		`json:"playerId"`
	PlayerSnake Snake 	`json:"playerSnake"`
	Foods []Food		`json:"foods"`
	OtherSnakes []Snake	`json:"otherSnakes"`
	Obstacles []Obstacle `json:"obstacles"`
}



func NewSnakeBoard() *SnakeBoard{
	snakeControllers := make(map[string]*SnakeController, 0)
	return &SnakeBoard{
		SnakeControllers: snakeControllers,
		Foods: make([]Food, 0),
		Obstacles: make([]Obstacle, 0),
		Width: 60,
		Height: 40,
		minimumFood: 4,
		numberOfFoodRange: 3,
	}
}

func(sb * SnakeBoard)AddPlayer(playerId string){
	sb.SnakeControllers[playerId] = NewSnakeController(NewSnake())
}

func (sb *SnakeBoard)GenerateFood(){
	snakes := make([]Snake, 0)
	for _, sc := range sb.SnakeControllers{
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
		if sb.isOccupied(newFood.Position, snakes, sb.Obstacles){
			continue
		}
		sb.Foods =append(sb.Foods, newFood)
	}
}

func (sb *SnakeBoard)isOccupied(p Point, snakes []Snake, obstacles []Obstacle) bool{
	for _, s := range(snakes){
		for _, body := range s.SnakeBody{
			if body.X == p.X && body.Y == p.Y {
				return true
			}
		}
		if s.SnakeHead.X == p.X && s.SnakeHead.Y == p.Y {
			return true
		}
		for _, obs := range obstacles {
			for _ , o := range obs.Object {
				if o.X == p.X && o.Y == p.Y {
					return true
				}
			}
		}
		
	}
	return false 
}


func (sb *SnakeBoard)ExecutePlayerMovement(playerId string, direction Direction){
	if sc, ok := sb.SnakeControllers[playerId]; ok {
		sc.KeyboardController(direction)
	}
}


func (sb *SnakeBoard)RunSnake(playerId string){
	sb.SnakeControllers[playerId].RunSnake(sb)
}

func (sb * SnakeBoard)GetSnakeBoard(playerId string) *SnakeBoardPlayerInformation{
	snakeController := sb.SnakeControllers[playerId]
	playerSnake := snakeController.Snake
	foods := sb.Foods
	obstacles := sb.Obstacles
	otherSnakes := make([]Snake, 0)
	for pId, sc := range(sb.SnakeControllers){
		if pId != playerId {
			otherSnakes = append(otherSnakes, *sc.Snake)
		} 
	}

	return &SnakeBoardPlayerInformation{
		PlayerId: playerId,
		PlayerSnake: *playerSnake,
		Foods: foods,
		Obstacles: obstacles,
		OtherSnakes: otherSnakes,
		
	}
	


}