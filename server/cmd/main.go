package main

import (
	"game-server/internal/api"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	PORT := ":8080"

	router.Use(cors.Default())

	// Register API routes
	api.PlayerRegisterRoutes(router)
	api.SnakeGameDataRoutes(router)
	api.MatchMakeRoutes(router)

	// Echo Server endpoint
	router.GET("/api", func(c *gin.Context) {
		c.String(http.StatusOK, "Echo from game-server")
	})

	// Start server
	log.Printf("Server Started at http://localhost%v", PORT)
	if err := router.Run(PORT); err != nil {
		log.Println("Server failure: ", err)
	}
}