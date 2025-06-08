package store

// import (
// 	"github.com/pramanandasarkar02/game-server/internal/game"
// 	"github.com/pramanandasarkar02/game-server/internal/models"
// )

// type PlayerStore interface {
// 	AddPlayer(player models.Player) error
// 	GetPlayer(id string) (models.Player, bool)
// 	GetAll() []models.Player
// }

// type GameStore interface {
// 	AddGame(g game.Game)
// 	GetGame(id string) (game.Game, bool)
// 	GetAll() []game.Game
// }

// type QueueStore interface {
// 	AddPlayer(gameID, playerID string)
// 	GetQueue(gameID string) []string
// 	RemovePlayers(gameID string, count int)
// 	GetAllQueues() map[string][]string
// }

// type MatchStore interface {
// 	AddMatch(match models.Match)
// 	GetMatch(id string) (models.Match, bool)
// 	GetMatchForPlayer(playerID string) (models.Match, bool)
// }