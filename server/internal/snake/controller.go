package snake

import "game-server/internal/service"


type SnakeController struct {
	Snake *Snake
	Player *service.Player 
}




type SnakeControllerResponse struct {
	Ok bool
	Msg string
}



func(sc *SnakeController) NewSnakeController(snake *Snake, player *service.Player) *SnakeController {
	return &SnakeController{
		Snake: snake,
		Player: player,
	}
}

func (sc *SnakeController)KeyboardController(option Direction)( error){
	sc.Snake.Controller(option)
	return nil
}