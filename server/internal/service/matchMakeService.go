package service

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var (
	GamePlayerRequirements = map[string]int{
		"single-snake-game":  1,
		"tic-tac-toe":        2,
		"snake": 2, 
		"four-snake-game":    4,
		"10-snake-game":      10,
		"random-snake-game":  0, // no support yet
	}
)

const (
	StatusQueued   = "queued"
	StatusInMatch  = "in_match"
	StatusIdle     = "idle"
)

type MatchMakeService struct {
	queue map[string]string // playerId -> gameId
	db    *sql.DB
	mu    sync.RWMutex
}

type GameEnv struct {
	GameId  string   `json:"gameId"`
	MatchId string   `json:"matchId"`
	Players []string `json:"players"`
}

type PlayerMatchResponse struct {
	PlayerId string  `json:"playerId"`
	GameEnv  GameEnv `json:"gameEnv"`
	Status   string  `json:"status"`
}

func NewMatchMakeService() *MatchMakeService {
	db, err := sql.Open("sqlite3", "./matches.db")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS matches (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			matchId TEXT UNIQUE,
			gameId TEXT,
			players TEXT
		);
		CREATE TABLE IF NOT EXISTS playerStatus (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			playerId TEXT UNIQUE,
			status TEXT,
			matchId TEXT
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	return &MatchMakeService{
		queue: make(map[string]string),
		db:    db,
	}
}

func (ms *MatchMakeService) AddQueue(playerId string, gameId string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// Check current player status
	status, matchId, err := ms.getPlayerStatus(playerId)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking player status: %v", err)
	}

	// If player is currently in a match, return error with existing match
	if status == StatusInMatch {
		return fmt.Errorf("player already in match: %v", matchId)
	}

	// Check if player is already in queue
	if _, exists := ms.queue[playerId]; exists {
		return fmt.Errorf("player already in queue with game id %v", ms.queue[playerId])
	}

	// Add to queue
	ms.queue[playerId] = gameId
	
	// Update player status to queued
	if err := ms.updatePlayerStatus(playerId, StatusQueued, ""); err != nil {
		delete(ms.queue, playerId)
		return fmt.Errorf("failed to update player status: %v", err)
	}

	log.Printf("Player %v added to queue for game %v", playerId, gameId)

	// Try to match
	ms.matchMake(gameId)
	return nil
}

func (ms *MatchMakeService) RemoveQueue(playerId string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.queue[playerId]; !exists {
		return fmt.Errorf("player %v not found in queue", playerId)
	}
	
	delete(ms.queue, playerId)
	
	// Update player status to idle
	if err := ms.updatePlayerStatus(playerId, StatusIdle, ""); err != nil {
		log.Printf("Failed to update player status to idle: %v", err)
	}

	log.Printf("%v removed from the queue", playerId)
	return nil
}

func (ms *MatchMakeService) GetMatch(playerId string) (*PlayerMatchResponse, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	// Check player status in DB
	status, matchId, err := ms.getPlayerStatus(playerId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("player %v has no active match", playerId)
		}
		return nil, fmt.Errorf("error getting player status: %v", err)
	}

	// If player is in match, load from DB
	if status == StatusInMatch && matchId != "" {
		gameEnv, err := ms.loadMatchFromDB(matchId)
		if err != nil {
			return nil, fmt.Errorf("failed to load match from DB: %v", err)
		}
		
		resp := &PlayerMatchResponse{
			PlayerId: playerId,
			GameEnv:  *gameEnv,
			Status:   StatusInMatch,
		}
		return resp, nil
	}

	return nil, fmt.Errorf("player %v not found in current matches", playerId)
}

func (ms *MatchMakeService) EndMatch(matchId string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// Load match from DB to get all players
	gameEnv, err := ms.loadMatchFromDB(matchId)
	if err != nil {
		return fmt.Errorf("failed to load match: %v", err)
	}

	// Update all players' status to idle
	for _, playerId := range gameEnv.Players {
		if err := ms.updatePlayerStatus(playerId, StatusIdle, ""); err != nil {
			log.Printf("Failed to update player %v status to idle: %v", playerId, err)
		}
	}

	log.Printf("Match %v ended, %d players freed", matchId, len(gameEnv.Players))
	return nil
}

func (ms *MatchMakeService) matchMake(gameId string) {
	var players []string
	for playerId, gId := range ms.queue {
		if gId == gameId {
			players = append(players, playerId)
		}
	}

	requiredPlayers := GamePlayerRequirements[gameId]

	// No support for random state player match yet
	if requiredPlayers == 0 {
		return
	}

	if len(players) >= requiredPlayers {
		// Take only the required number of players
		selectedPlayers := players[:requiredPlayers]
		
		matchId := fmt.Sprintf("match-%v", uuid.New())
		log.Printf("Creating match %v for game %v with players: %v", matchId, gameId, selectedPlayers)

		gameEnv := GameEnv{
			GameId:  gameId,
			MatchId: matchId,
			Players: selectedPlayers,
		}

		if err := ms.saveMatchToDB(gameEnv); err != nil {
			log.Printf("Error saving match to DB: %v", err)
			return
		}

		// Remove players from queue and update their status
		for _, p := range selectedPlayers {
			delete(ms.queue, p)
			
			// Update player status to in_match
			if err := ms.updatePlayerStatus(p, StatusInMatch, matchId); err != nil {
				log.Printf("Failed to update player %v status: %v", p, err)
			}
		}
		
		log.Printf("Match %v created successfully", matchId)
	}
}

func (ms *MatchMakeService) saveMatchToDB(env GameEnv) error {
	playerList := strings.Join(env.Players, ",")
	_, err := ms.db.Exec(`
		INSERT INTO matches (matchId, gameId, players)
		VALUES (?, ?, ?)
	`, env.MatchId, env.GameId, playerList)

	if err != nil {
		return fmt.Errorf("failed to insert match: %v", err)
	}
	return nil
}

func (ms *MatchMakeService) loadMatchFromDB(matchId string) (*GameEnv, error) {
	var gameId, playerList string
	err := ms.db.QueryRow(`
		SELECT gameId, players FROM matches WHERE matchId = ?
	`, matchId).Scan(&gameId, &playerList)

	if err != nil {
		return nil, fmt.Errorf("failed to load match: %v", err)
	}

	players := strings.Split(playerList, ",")
	return &GameEnv{
		GameId:  gameId,
		MatchId: matchId,
		Players: players,
	}, nil
}

func (ms *MatchMakeService) getPlayerStatus(playerId string) (string, string, error) {
	var status, matchId string
	err := ms.db.QueryRow(`
		SELECT status, matchId FROM playerStatus WHERE playerId = ?
	`, playerId).Scan(&status, &matchId)

	return status, matchId, err
}

func (ms *MatchMakeService) updatePlayerStatus(playerId, status, matchId string) error {
	_, err := ms.db.Exec(`
		INSERT INTO playerStatus (playerId, status, matchId)
		VALUES (?, ?, ?)
		ON CONFLICT(playerId) DO UPDATE SET
			status = excluded.status,
			matchId = excluded.matchId
	`, playerId, status, matchId)

	if err != nil {
		return fmt.Errorf("failed to update player status: %v", err)
	}
	return nil
}

func (ms *MatchMakeService) Close() error {
	if ms.db != nil {
		return ms.db.Close()
	}
	return nil
}