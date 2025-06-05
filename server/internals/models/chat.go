package models


type ChatMessage struct {
	MatchID string `json:"matchId"`
	PlayerID string `json:"playerId"`
	Content string `json:"content"`
	Timestamp string `json:"timestamp"`
}