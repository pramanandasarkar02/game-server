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
	minimumFood int 
	numberOfFoodRange int 
	obstacles []Obstacle
}


func NewSnakeBoard() *SnakeBoard{
	return &SnakeBoard{
		
	}
}

func (sb *SnakeBoard)generateFood(){
	
}