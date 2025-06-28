package models
type ChatMessage struct{
	MatchID string `json:"matchId"`
	PlayerID string `josn:"playerId"`
	Content string `josn:"content"`
	Timestamp string `json:"timestamp"`
}

