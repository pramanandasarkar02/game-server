package services

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pramanandasarkar02/game-server/internal/dtos"
	"github.com/pramanandasarkar02/game-server/internal/store"
	"golang.org/x/crypto/bcrypt"
)


type PlayerService struct {
	store *store.PlayerStore
}


func NewPlayerService(store *store.PlayerStore) *PlayerService {
	return &PlayerService{
		store: store,
	}
}

func (s *PlayerService) ConnectPlayer(req dtos.PlayerConnectionRequest) (dtos.PlayerConnectionResponse, error) {
	
	if err := req.Validate(); err != nil {
		return dtos.PlayerConnectionResponse{}, err
	}

	// print all players
	players := s.store.GetAll()
	fmt.Println("Players:", players)

	player, err := s.store.GetPlayerByUsername(req.Username)

	if err != nil {
		return dtos.PlayerConnectionResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(req.Password)); err != nil {
		return dtos.PlayerConnectionResponse{}, fmt.Errorf("invalid credentials")
	}

	token, err := s.generateToken(player.UserId, req.Username)
	if err != nil {
		return dtos.PlayerConnectionResponse{}, err
	}
	if err := s.store.SaveSession(player.UserId, token); err != nil {
		return dtos.PlayerConnectionResponse{}, fmt.Errorf("failed to save session: %v", err)
	}

	return dtos.PlayerConnectionResponse{
		Token:     token,
		UserId:    player.UserId,
		Username:  player.Username,
	}, nil
}

func (s *PlayerService) RegisterPlayer(req dtos.PlayerRegisterRequest) (dtos.PlayerRegisterResponse, error) {
	if err := req.Validate(); err != nil {
		return dtos.PlayerRegisterResponse{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dtos.PlayerRegisterResponse{}, fmt.Errorf("failed to hash password: %v", err)
	}

	playerStore := dtos.PlayerRegisterStore{
		Username: req.Username,
		Password: string(hashedPassword),
	}
	if err := s.store.CreateNewPlayer(playerStore); err != nil {
		return dtos.PlayerRegisterResponse{}, fmt.Errorf("failed to create player: %v", err)
	}
	
	player, err := s.store.GetPlayerByUsername(req.Username)
	
	if err != nil {
		return dtos.PlayerRegisterResponse{}, err
	}

	token, err := s.generateToken(player.UserId, player.Username)
	if err != nil {
		return dtos.PlayerRegisterResponse{}, err
	}
	return dtos.PlayerRegisterResponse{
		Token:     token,
		UserId:    playerStore.Username,
		Username:  playerStore.Username,
	}, nil

}
















// -------------------------------------------------------------------------- //

func (s *PlayerService) generateToken(userId, username string) (string, error) {
	// Generate JWT token
	token, err := generateJWT(userId, username)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}
	return token, nil
}

func generateJWT(userId, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userId,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte("SECRET_KEY")) 
}