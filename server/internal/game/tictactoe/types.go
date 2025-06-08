package tictactoe

type TicTacToeState struct {
	Board   [9]string `json:"board"`
	Turn    string    `json:"turn"`
	Winner  string    `json:"winner"`
	IsDraw  bool      `json:"isDraw"`
	Players []string  `json:"players"`
}