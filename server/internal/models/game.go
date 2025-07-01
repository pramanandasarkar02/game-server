package models

import "github.com/google/uuid"

type Game struct {
    ID             string `json:"id"`
    Title          string `json:"title"`
    RequiredPlayer int    `json:"requiredPlayer"`
    MetaData       string `json:"metaData"`
}

func NewGame(title string, requiredPlayer int, metadata string) *Game {
    return &Game{
        ID:             "game-" + uuid.New().String(),
        Title:          title,
        RequiredPlayer: requiredPlayer,
        MetaData:       metadata,
    }
}