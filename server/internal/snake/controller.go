package snake

import "game-server/internal/service"


type SnakeController struct {
	Snake *Snake
	Player *service.Player 
}

type ControllerOption string


type SnakeControllerResponse struct {
	Ok bool
	Msg string
}


const (
	LEFT_CONTROLLER ControllerOption = "left"
	RIGHT_CONTROLLER ControllerOption = "right"
	UP_CONTROLLER ControllerOption = "top"
	DOWN_CONTROLLER ControllerOption = "down"

)

func(sc *SnakeController) NewSnakeController(snake *Snake, player *service.Player) *SnakeController {
	return &SnakeController{
		Snake: snake,
		Player: player,
	}
}

func (sc *SnakeController)KeyboardController(option ControllerOption)( error){
	sc.Snake.Controller(option)
	return nil
}