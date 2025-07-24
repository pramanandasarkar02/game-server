package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pramanandasarkar02/game-server/internal/dtos"
	"github.com/pramanandasarkar02/game-server/internal/services"

	// "github.com/pramanandasarkar02/game-server/internal/models"
	// "github.com/pramanandasarkar02/game-server/internal/store"
	"github.com/pramanandasarkar02/game-server/pkg/logger"
)


func PlayerRegisterHandler(service  *services.PlayerService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var playerRegistrationRequest dtos.PlayerRegisterRequest
        if err := c.ShouldBindJSON(&playerRegistrationRequest); err != nil {
            logger.Error("Invalid player registration data: %v", err)
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }

        // Validate request
        if err := playerRegistrationRequest.Validate(); err != nil {
            logger.Error("Validation failed: %v", err)
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }

        playerResponse, err := service.RegisterPlayer(&playerRegistrationRequest)
        if err != nil {
            logger.Error("Failed to register player: %v", err)
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }

        logger.Info("Player %s registered", playerRegistrationRequest.Username)
        c.JSON(200, gin.H{"message": "Player registered successfully", "player": playerResponse})
    }
}

func PlayerConnectionHandler(service  *services.PlayerService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var playerConnectionRequest dtos.PlayerConnectionRequest
        if err := c.ShouldBindJSON(&playerConnectionRequest); err != nil {
            logger.Error("Invalid player connection data: %v", err)
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }

        // Validate request
        if err := playerConnectionRequest.Validate(); err != nil {
            logger.Error("Validation failed: %v", err)
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }

        playerResponse, err := service.ConnectPlayer(&playerConnectionRequest)
        if err != nil {
            logger.Error("Failed to connect player: %v", err)
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }

        logger.Info("Player %s connected", playerConnectionRequest.Username)
        c.JSON(200, gin.H{"message": "Player connected successfully", "player": playerResponse})
    }
}

func PlayerAuthValidationHandler(service  *services.PlayerService) gin.HandlerFunc {
    return func(c *gin.Context) {
        var playerAuthValidationRequest dtos.PlayerAuthValidationRequest
        if err := c.ShouldBindJSON(&playerAuthValidationRequest); err != nil {
            logger.Error("Invalid player authentication data: %v", err)
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        // Validate request
        if err := playerAuthValidationRequest.Validate(); err != nil {
            logger.Error("Validation failed: %v", err)
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        playerResponse, err := service.ValidateToken(playerAuthValidationRequest.Token)
        if err != nil {
            logger.Error("Failed to validate player: %v", err)
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        logger.Info("Player %s authenticated", playerResponse.Username)
        c.JSON(200, gin.H{"message": "Player authenticated successfully", "player": playerResponse})
    }
}

func PlayerDisconnectionHandler(service  *services.PlayerService) gin.HandlerFunc {
    return func(c *gin.Context) {
        playerId := c.Param("playerId")
        err := service.DisconnectPlayer(playerId)
        if err != nil {
            logger.Error("Failed to disconnect player: %v", err)
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        logger.Info("Player %s disconnected", playerId)
        c.JSON(200, gin.H{"message": "Player disconnected successfully"})
    }
}


func PlayerInfoHandler(service  *services.PlayerService) gin.HandlerFunc {
    return func(c *gin.Context) {
        playerId := c.Param("playerId")
        fmt.Printf("playerId: %s\n", playerId)
        player, err := service.GetPlayerInfo(playerId)
        if err != nil {
            logger.Error("Failed to get player: %v", err)
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, gin.H{"player": player})
    }
}
