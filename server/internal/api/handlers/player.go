package handlers

import (
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




// func PlayerHandler(service  *services.PlayerService) gin.HandlerFunc {
//     return func(c *gin.Context) {
//         var playerConnectionRequest dtos.PlayerConnectionRequest
//         if err := c.ShouldBindJSON(&playerConnectionRequest); err != nil {
//             logger.Error("Invalid player connection data: %v", err)
//             c.JSON(400, gin.H{"error": err.Error()})
//             return
//         }

//         // Validate request
//         if err := playerConnectionRequest.Validate(); err != nil {
//             logger.Error("Validation failed: %v", err)
//             c.JSON(400, gin.H{"error": err.Error()})
//             return
//         }

//         playerResponse, err := service.ConnectPlayer(playerConnectionRequest)
//         if err != nil {
//             logger.Error("Failed to connect player: %v", err)
//             c.JSON(400, gin.H{"error": err.Error()})
//             return
//         }

//         logger.Info("Player %s connected", playerConnectionRequest.Username)
//         c.JSON(200, gin.H{"message": "Player connected successfully", "player": playerResponse})
//     }
// }


// func PlayerConnectionHandler(service  *services.PlayerService) gin.HandlerFunc {
//     return func(c *gin.Context) {
//         var playerConnectionRequest dtos.PlayerConnectionRequest
//         if err := c.ShouldBindJSON(&playerConnectionRequest); err != nil {
//             logger.Error("Invalid player connection data: %v", err)
//             c.JSON(400, gin.H{"error": err.Error()})
//             return
//         }

//         // Validate request
//         if err := playerConnectionRequest.Validate(); err != nil {
//             logger.Error("Validation failed: %v", err)
//             c.JSON(400, gin.H{"error": err.Error()})
//             return
//         }

//         playerResponse, err := service.ConnectPlayer(playerConnectionRequest)
//         if err != nil {
//             logger.Error("Failed to connect player: %v", err)
//             c.JSON(400, gin.H{"error": err.Error()})
//             return
//         }

//         logger.Info("Player %s connected", playerConnectionRequest.Username)
//         c.JSON(200, gin.H{"message": "Player connected successfully", "player": playerResponse})
//     }
// }

// func PlayerRegisterHandler(service  *services.PlayerService) gin.HandlerFunc {
//     return func(c *gin.Context) {
//         var playerRegisterRequest dtos.PlayerRegisterRequest
//         if err := c.ShouldBindJSON(&playerRegisterRequest); err != nil {
//             logger.Error("Invalid player connection data: %v", err)
//             c.JSON(400, gin.H{"error": err.Error()})
//             return
//         }

//         // Validate request
//         if err := playerRegisterRequest.Validate(); err != nil {
//             logger.Error("Validation failed: %v", err)
//             c.JSON(400, gin.H{"error": err.Error()})
//             return
//         }

//         playerResponse, err := service.RegisterPlayer(playerRegisterRequest)
//         if err != nil {
//             logger.Error("Failed to connect player: %v", err)
//             c.JSON(400, gin.H{"error": err.Error()})
//             return
//         }

//         logger.Info("Player %s connected", playerResponse.Username)
//         c.JSON(200, gin.H{"message": "Player connected successfully", "player": playerResponse})
//     }
// }












// --------------------------------------------------------------------------- //


// func GetPlayersHandler(store *store.PlayerStore) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		players := store.GetAll()
// 		c.JSON(200, gin.H{"players": players})
// 	}
// }