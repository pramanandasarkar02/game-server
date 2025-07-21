package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pramanandasarkar02/game-server/config"
	"github.com/pramanandasarkar02/game-server/internal/api/handlers"
	"github.com/pramanandasarkar02/game-server/internal/services"
	"github.com/pramanandasarkar02/game-server/internal/store"
	"github.com/pramanandasarkar02/game-server/internal/websocket"
)

func NewRouter(playerService *services.PlayerService, playerStore *store.PlayerStore, gameStore *store.GameStore, queueStore *store.QueueStore, matchStore *store.MatchStore, cfg *config.Config) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	wsManager := websocket.NewWebSocketManager(cfg, gameStore, matchStore)

	router.GET("/ping", handlers.PingHandler())
	router.POST("/register", handlers.PlayerRegisterHandler(playerService))
	// router.POST("/connect", handlers.PlayerConnectionHandler(playerService))
	// router.GET("/players", handlers.GetPlayersHandler(playerStore))
	router.GET("/games", handlers.GetGamesHandler(gameStore))
	router.GET("/queue", handlers.GetQueueHandler(queueStore))
	router.POST("/queue/join", handlers.JoinQueueHandler(playerStore, gameStore, queueStore))
	router.GET("/match/:userId", handlers.GetMatchForUserHandler(matchStore))
	router.GET("/running-match/:matchId/:userId", handlers.GetMatchHandler(matchStore))
	router.GET("/chat/:matchId/:playerId", wsManager.HandleConnection)

	return router
}