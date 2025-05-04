package model


type Game struct {
	Id string `json:"id"`
	Board [3][3]string `json:"board"`
	CurrentPlayer string `json:"currentPlayer"`
	Status string `json:"status"`
	Winner string `json:"winner"`
}