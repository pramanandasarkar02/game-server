package main

import (
	"github.com/pramanandasarkar02/game-server/config"
	"github.com/pramanandasarkar02/game-server/internal/api"
	"github.com/pramanandasarkar02/game-server/internal/game"
	"github.com/pramanandasarkar02/game-server/internal/matchmaking"
	"github.com/pramanandasarkar02/game-server/internal/models"
	"github.com/pramanandasarkar02/game-server/internal/services"
	"github.com/pramanandasarkar02/game-server/internal/store"
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
	playerStore := store.NewPlayerStore()
	gameStore := store.NewGameStore()
	queueStore := store.NewQueueStore()
	matchStore := store.NewMatchStore()


	// create demo user 
	playerStore.AddPlayer(models.Player{
		Username: "admin",
		Password: "password",
	})



	// initialize services
	playerService := services.NewPlayerService(playerStore)

	// Initialize games
	gameStore.AddGame(game.TicTacToeGame())

	// Start matchmaker
	go matchmaking.StartMatchmaker(gameStore, queueStore, matchStore)

	// Start API server
	router := api.NewRouter(playerService, playerStore, gameStore, queueStore, matchStore, cfg)



	if err := router.Run(":" + cfg.Port); err != nil {
		logger.Fatal("Failed to run server: %v", err)
	}
}