package handler

import (
	"game-server/internal/service"
	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	playerService *service.PlayerService
}

func NewPlayerHandler(ps *service.PlayerService) *PlayerHandler {
	return &PlayerHandler{
		playerService: ps,
	}
}

// Login handler
func (ph *PlayerHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	player, err := ph.playerService.Login(req)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, player)
}

// Logout handler
func (ph *PlayerHandler) Logout(c *gin.Context) {
	var req service.LogOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	message, err := ph.playerService.Logout(req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": message})
}

// Signup handler
func (ph *PlayerHandler) SignUp(c *gin.Context) {
	var req service.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	player, err := ph.playerService.Signup(req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, player)
}