package main

import (
	"github.com/pramanandasarkar02/game-server/config"
	"github.com/pramanandasarkar02/game-server/internal/api"
	"github.com/pramanandasarkar02/game-server/internal/game"
	"github.com/pramanandasarkar02/game-server/internal/matchmaking"
	"github.com/pramanandasarkar02/game-server/internal/store/memory"
	"github.com/pramanandasarkar02/game-server/pkg/logger"
)

func main() {
	// Initialize logger
	logger.Init()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config: %v", err)
	}

	// Initialize stores
	playerStore := memory.NewPlayerStore()
	gameStore := memory.NewGameStore()
	queueStore := memory.NewQueueStore()
	matchStore := memory.NewMatchStore()

	// Initialize games
	gameStore.AddGame(game.TicTacToeGame())

	// Start matchmaker
	go matchmaking.StartMatchmaker(gameStore, queueStore, matchStore)

	// Start API server
	router := api.NewRouter(playerStore, gameStore, queueStore, matchStore, cfg)
	if err := router.Run(":" + cfg.Port); err != nil {
		logger.Fatal("Failed to run server: %v", err)
	}
}