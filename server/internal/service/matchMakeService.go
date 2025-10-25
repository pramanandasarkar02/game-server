package service

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var (
	MINIMUM_MATCH_PLAYER = 2
)

type MatchMakeService struct {
	queue   map[string]string // playerId -> gameId
	matches map[string]GameEnv
	db      *sql.DB
}

type GameEnv struct {
	GameId  string   `json:"gameId"`
	MatchId string   `json:"matchId"`
	Players []string `json:"players"`
}

type PlayerMatchResponse struct {
	PlayerId string  `json:"playerId"`
	GameEnv  GameEnv `json:"gameEnv"`
}

func NewMatchMakeService() *MatchMakeService {
	db, err := sql.Open("sqlite3", "./matches.db")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS matches (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			match_id TEXT UNIQUE,
			game_id TEXT,
			players TEXT
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	return &MatchMakeService{
		queue:   make(map[string]string),
		matches: make(map[string]GameEnv),
		db:      db,
	}
}

func (ms *MatchMakeService) AddQueue(playerId string, gameId string) error {
	if _, exists := ms.queue[playerId]; exists {
		return fmt.Errorf("player already in queue with game id %v", ms.queue[playerId])
	}

	ms.queue[playerId] = gameId
	log.Printf("Player %v added to queue for game %v", playerId, gameId)

	ms.matchMake(gameId)
	return nil
}

func (ms *MatchMakeService) RemoveQueue(playerId string) error {
	if _, exists := ms.queue[playerId]; !exists {
		return fmt.Errorf("player %v not found in queue", playerId)
	}
	delete(ms.queue, playerId)

	log.Printf("%v removed from the queue", playerId)
	return nil
}

func (ms *MatchMakeService) GetMatch(playerId string) (*PlayerMatchResponse, error) {
	gameEnv, exists := ms.matches[playerId]
	if !exists {
		return nil, fmt.Errorf("playerId %v not found in current matches", playerId)
	}

	resp := &PlayerMatchResponse{
		PlayerId: playerId,
		GameEnv:  gameEnv,
	}
	return resp, nil
}

func (ms *MatchMakeService) matchMake(gameId string) {
	var players []string
	for playerId, gId := range ms.queue {
		if gId == gameId {
			players = append(players, playerId)
		}
	}

	if len(players) >= MINIMUM_MATCH_PLAYER {
		matchId := fmt.Sprintf("match-%v", uuid.New())
		log.Printf("Creating match %v for game %v with players: %v", matchId, gameId, players)

		gameEnv := GameEnv{
			GameId:  gameId,
			MatchId: matchId,
			Players: players,
		}

		if err := ms.saveMatchToDB(gameEnv); err != nil {
			log.Printf("Error saving match to DB: %v", err)
			return
		}

		for _, p := range players {
			delete(ms.queue, p)
			ms.matches[p] = gameEnv
		}
		log.Printf("Match %v created successfully", matchId)
	}
}

func (ms *MatchMakeService) saveMatchToDB(env GameEnv) error {
	playerList := strings.Join(env.Players, ",")
	_, err := ms.db.Exec(`
		INSERT INTO matches (match_id, game_id, players)
		VALUES (?, ?, ?)
	`, env.MatchId, env.GameId, playerList)

	if err != nil {
		return fmt.Errorf("failed to insert match: %v", err)
	}
	return nil
}
