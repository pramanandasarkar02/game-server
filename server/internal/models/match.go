package models

type Match struct {
	ID      string   `json:"id"`
	GameID  string   `json:"gameId"`
	Players []string `json:"players"`
}