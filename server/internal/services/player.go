package services

import (
	// "fmt"
	// "time"
	// "github.com/dgrijalva/jwt-go"
	"github.com/pramanandasarkar02/game-server/internal/dtos"
	"github.com/pramanandasarkar02/game-server/internal/models"
	"github.com/pramanandasarkar02/game-server/internal/store"
	// "golang.org/x/crypto/bcrypt"
)

// player service to handle all player related operations
type PlayerService struct {
	store *store.PlayerStore
}

// create new instance of player service
func NewPlayerService(store *store.PlayerStore) *PlayerService {
	return &PlayerService{
		store: store,
	}
}

// register player with username and password
func (ps *PlayerService) RegisterPlayer(playerRequest *dtos.PlayerRegisterRequest) (dtos.PlayerRegisterResponse, error) {

	if err := playerRequest.Validate(); err != nil {
		return dtos.PlayerRegisterResponse{}, err
	}

	player := *models.NewPlayer(playerRequest.Username, playerRequest.Password)

	if err := ps.store.AddPlayer(player); err != nil {
		return dtos.PlayerRegisterResponse{}, err
	}
	
	return dtos.PlayerRegisterResponse{
		UserId: player.ID,
		Username: player.Username,
	}, nil
	
}

// connect player with username and password
func (ps *PlayerService) ConnectPlayer(playerRequest *dtos.PlayerConnectionRequest) (dtos.PlayerConnectionResponse, error) {
	if err := playerRequest.Validate(); err != nil {
		return dtos.PlayerConnectionResponse{}, err
	}
	playerId, err := ps.store.GetPlayerIdByUsername(playerRequest.Username)
	if err != nil {
		return dtos.PlayerConnectionResponse{}, err
	}
	player, err := ps.store.GetPlayerById(playerId)
	if err != nil {
		return dtos.PlayerConnectionResponse{}, err
	}
	return dtos.PlayerConnectionResponse{
		UserId: player.ID,
		Username: player.Username,
	}, nil
}

