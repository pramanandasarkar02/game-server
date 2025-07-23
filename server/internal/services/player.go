package services

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pramanandasarkar02/game-server/internal/dtos"
	"github.com/pramanandasarkar02/game-server/internal/models"
	"github.com/pramanandasarkar02/game-server/internal/store"
	"golang.org/x/crypto/bcrypt"

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


// generateJWT token for a player
func generateToken(userId, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"username": username,
		"exp":      jwt.TimeFunc().Add(time.Hour * 24).Unix(),
		"iat":      jwt.TimeFunc().Unix(),
	})
	return token.SignedString([]byte("secret"))
}


// register player with username and password
func (ps *PlayerService) RegisterPlayer(playerRequest *dtos.PlayerRegisterRequest) (dtos.PlayerRegisterResponse, error) {

	if err := playerRequest.Validate(); err != nil {
		return dtos.PlayerRegisterResponse{}, err
	}

	playerId, _ := ps.store.GetPlayerIdByUsername(playerRequest.Username)
	if playerId != "" {
		return dtos.PlayerRegisterResponse{}, fmt.Errorf("player with username %s already exists", playerRequest.Username)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(playerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return dtos.PlayerRegisterResponse{}, err
	}

	player := *models.NewPlayer(playerRequest.Username, string(hashedPassword))

	if err := ps.store.AddPlayer(player); err != nil {
		return dtos.PlayerRegisterResponse{}, err
	}
	// generate token save session info in db
	token, err := generateToken(player.ID, player.Username)
	if err != nil {
		return dtos.PlayerRegisterResponse{}, err
	}
	
	return dtos.PlayerRegisterResponse{
		UserId: player.ID,
		Username: player.Username,
		Token: token,
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
	if err := bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(playerRequest.Password)); err != nil {
		return dtos.PlayerConnectionResponse{}, err
	}

	// generate token save session info in db
	token, err := generateToken(player.ID, player.Username)
	if err != nil {
		return dtos.PlayerConnectionResponse{}, err
	}

	ps.store.SaveToken(playerId, token)

	return dtos.PlayerConnectionResponse{
		UserId: player.ID,
		Username: player.Username,
		Token: token,
	}, nil
}

func (ps *PlayerService) DisconnectPlayer(playerId string) error {
	return ps.store.DeleteToken(playerId)
}

func (ps *PlayerService) ValidateToken(tokenString string) (*dtos.PlayerConnectionResponse, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &dtos.PlayerConnectionResponse{
			UserId: claims["userId"].(string),
			Username: claims["username"].(string),
			Token: tokenString,
		}, nil
	}
	return nil, err
}