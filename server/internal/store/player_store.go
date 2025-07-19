package store

import (
	"errors"
	"log"
	"sync"

	"github.com/pramanandasarkar02/game-server/internal/dtos"
	"github.com/pramanandasarkar02/game-server/internal/models"
)

// players data save in memory
type PlayerStore struct {
	players []models.Player
	mutex   sync.RWMutex
	activeSessions map[string]string
}



func NewPlayerStore() *PlayerStore {
	return &PlayerStore{
		players: make([]models.Player, 0),
		activeSessions: make(map[string]string),
	}
}


func (s *PlayerStore) CreatePlayer(playerDto dtos.CreatePlayerDto) {
	if err := playerDto.Validate(); err != nil {
		return 
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	// save in memory and postgress db


}

func (s* PlayerStore) GetPlayerProfileInfo(playerProfileInfoDto dtos.PlayerProfileInfoDto) {
	
}

func (s *PlayerStore) AddPlayer(player models.Player) error {
	if err := player.Validate(); err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.players = append(s.players, player)
	log.Printf("Player added: %+v, Total players: %d", player, len(s.players))
	return nil
}

func (s *PlayerStore) GetPlayer(id string) (models.Player, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, p := range s.players {
		if p.ID == id {
			return p, true
		}
	}
	log.Printf("Player with ID %s not found", id)
	return models.Player{}, false
}

func (s *PlayerStore) GetAll() []models.Player {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	playersCopy := make([]models.Player, len(s.players))
	copy(playersCopy, s.players)
	log.Printf("Returning %d players: %+v", len(playersCopy), playersCopy)
	return playersCopy // Return the copy to prevent external modification
}

// func (s *PlayerStore) ConnectPlayer(playerConnectionRequest dtos.PlayerConnectionRequest) (dtos.PlayerConnectionResponse, error) {
//     s.mutex.Lock()
//     defer s.mutex.Unlock()

//     // Validate credentials against PostgreSQL
//     var storedPlayer models.Player // Assume models.Player exists
//     err := s.db.Where("username = ?", playerConnectionRequest.Username).First(&storedPlayer).Error
//     if err != nil {
//         if errors.Is(err, gorm.ErrRecordNotFound) {
//             return dtos.PlayerConnectionResponse{}, fmt.Errorf("invalid credentials")
//         }
//         return dtos.PlayerConnectionResponse{}, fmt.Errorf("database error: %v", err)
//     }

//     // Verify password (assuming password is hashed in DB)
//     if err := bcrypt.CompareHashAndPassword([]byte(storedPlayer.Password), []byte(playerConnectionRequest.Password)); err != nil {
//         return dtos.PlayerConnectionResponse{}, fmt.Errorf("invalid credentials")
//     }

//     // Generate JWT token
//     token, err := generateJWT(storedPlayer.ID, playerConnectionRequest.Username)
//     if err != nil {
//         return dtos.PlayerConnectionResponse{}, fmt.Errorf("failed to generate token: %v", err)
//     }

//     // Store session in Redis
//     err = s.redis.Set(context.Background(), fmt.Sprintf("session:%s", storedPlayer.ID), token, 24*time.Hour).Err()
//     if err != nil {
//         return dtos.PlayerConnectionResponse{}, fmt.Errorf("failed to store session: %v", err)
//     }

//     return dtos.PlayerConnectionResponse{
//         Token:    token,
//         UserId:   storedPlayer.ID,
//         Username: playerConnectionRequest.Username,
//     }, nil
// }

// // Helper function to generate JWT
// func generateJWT(userID, username string) (string, error) {
//     token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//         "user_id":  userID,
//         "username": username,
//         "exp":      time.Now().Add(24 * time.Hour).Unix(),
//     })

//     return token.SignedString([]byte("your-secret-key")) // Replace with actual secret
// }



func (s *PlayerStore)CreateNewPlayer(player dtos.PlayerRegisterStore)  error {
	// Validate player data
	if err := player.Validate(); err != nil {
		return err
	}

	// Save player data to database
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.players = append(s.players, models.Player{
		ID:       player.Username,
		Name:     player.Username,
		Password: player.Password,
	})
	return nil
}

func (s *PlayerStore)GetPlayerByUsername(username string) (dtos.PlayerConnectionResponse, error) {
	for _, player := range s.players {
		if player.Name == username {
			return dtos.PlayerConnectionResponse{
				Token:    "token",
				UserId:   player.ID,
				Username: player.Name,
			}, nil
		}
	}
	return dtos.PlayerConnectionResponse{}, errors.New("player not found")
}


func (s *PlayerStore) SaveSession(userId, token string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.activeSessions[userId] = token
	return nil
}

