package snake



type SnakeController struct {
	Snake *Snake

}




type SnakeControllerResponse struct {
	Ok bool
	Msg string
}



func NewSnakeController(snake *Snake) *SnakeController {
	return &SnakeController{
		Snake: snake,
	}
}
func (sc* SnakeController)RunSnake(snakeBoard *SnakeBoard)(bool, string){
	isCol, msg := sc.Snake.Movement(snakeBoard)
	return isCol, msg
}

func (sc *SnakeController)KeyboardController(option Direction)( error){
	sc.Snake.Controller(option)
	return nil
}
