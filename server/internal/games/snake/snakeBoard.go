package snake

import (
)

type Food struct{
	Position Point 
	Value int
}

type Obstacle struct{
	Object []Point
}


type SnakeBoard struct{
	SnakeControllers []SnakeController
	Foods []Food
	Obstacle []Point
	Width int 
	Height int
}

func generateFood()