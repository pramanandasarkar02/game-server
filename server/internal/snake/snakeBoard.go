package snake

import "math/rand/v2"

type Food struct{
	Position Point 
	Value int
}

type Obstacle struct{
	Object []Point
}


type SnakeBoard struct{
	SnakeControllers map[string]SnakeController
	Foods []Food
	Width int 
	Height int
	minimumFood int 
	numberOfFoodRange int 
	Obstacles []Obstacle
}

type SnakeBoardPlayerInformation struct{
	PlayerSnake Snake
	PlayerId string
	Foods []Food
	OtherSnake []Snake
	Obstacles []Obstacle 
}



func NewSnakeBoard() *SnakeBoard{


	return &SnakeBoard{
		
	}
}

func (sb *SnakeBoard)generateFood(){
	snakes := make([]Snake, 0)
	for _, sc := range sb.SnakeControllers{
		snakes = append(snakes, *sc.Snake)
	}
	numberOfFood := sb.minimumFood + rand.IntN(sb.numberOfFoodRange)
	for _ := range(numberOfFood){
		
	}



}


func (sb *SnakeBoard)ExecutePlayerMovement(playerId string, direction Direction){
	snakeController := sb.SnakeControllers[playerId]
	snakeController.KeyboardController(direction)
}

func (sb * SnakeBoard)GetSnakeBoard(playerId string) SnakeBoardPlayerInformation{
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

	snakeBoardPlayerInformation := SnakeBoardPlayerInformation{
		PlayerId: playerId,
		PlayerSnake: *playerSnake,
		Foods: foods,
		Obstacles: obstacles,
		OtherSnake: otherSnakes,
	}
	return snakeBoardPlayerInformation


}