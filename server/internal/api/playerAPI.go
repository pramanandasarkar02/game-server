package api

import (
	"game-server/internal/handler"
	"game-server/internal/service"

	"github.com/gin-gonic/gin"
)

func PlayerRegisterRoutes(router *gin.Engine) {
	// create service
	playerService := service.NewPlayerService()

	// inject service into handler
	playerHandler := handler.NewPlayerHandler(playerService)

	// register routes
	router.POST("/api/login", playerHandler.Login)
	router.POST("/api/logout", playerHandler.Logout)
	router.POST("/api/signup", playerHandler.SignUp)
	
}
