package dtos

import (
	"fmt"
	"regexp"
	"strings"
)

type CreatePlayerDto struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
}

func (dto *CreatePlayerDto) Validate() error {
	if strings.TrimSpace(dto.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if len(dto.Name) < 3 || len(dto.Name) > 50 {
		return fmt.Errorf("name must be between 3 and 50 characters")
	}
	if len(dto.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	if !isValidEmail(dto.Email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

type PlayerProfileInfoDto struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Level        float32  `json:"level"`
	MatchHistory []string `json:"matchHistory"`
	PlayerStatus string   `json:"playerStatus"`
	WinRate      float32  `json:"winRate"`
	MatchCount   int      `json:"matchCount"`
	CreatedAt    string   `json:"createdAt"`
	UpdatedAt    string   `json:"updatedAt"`
}

type PlayerMatchInfoDto struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Level        float32         `json:"level"`
	State        string          `json:"state"`
	MatchHistory map[string]bool `json:"matchHistory"`
	WinRate      float32         `json:"winRate"`
	MatchCount   int             `json:"matchCount"`
}

type PlayerAuthUpdateDto struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=6"`
	Email       string `json:"email,omitempty" validate:"omitempty,email"`
}

func (dto *PlayerAuthUpdateDto) Validate() error {
	if dto.OldPassword == "" {
		return fmt.Errorf("old password is required")
	}
	if len(dto.NewPassword) < 6 {
		return fmt.Errorf("new password must be at least 6 characters")
	}
	if dto.Email != "" && !isValidEmail(dto.Email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

type PlayerMatchUpdateDto struct {
	MatchID string `json:"matchId" validate:"required"`
	Won     bool   `json:"won"`
}

func (dto *PlayerMatchUpdateDto) Validate() error {
	if strings.TrimSpace(dto.MatchID) == "" {
		return fmt.Errorf("match ID is required")
	}
	return nil
}

type PlayerLevelUpdateDto struct {
	Level float32 `json:"level" validate:"required,min=1"`
}

func (dto *PlayerLevelUpdateDto) Validate() error {
	if dto.Level < 1 {
		return fmt.Errorf("level must be at least 1")
	}
	return nil
}

type PlayerStateUpdateDto struct {
	State string `json:"state" validate:"required"`
}

func (dto *PlayerStateUpdateDto) Validate() error {
	validStates := []string{"InGame", "InQuery", "Offline", "Online"}
	for _, state := range validStates {
		if dto.State == state {
			return nil
		}
	}
	return fmt.Errorf("invalid state: %s, valid states are: %v", dto.State, validStates)
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}



type PlayerConnectionRequest struct {
	Username     string `json:"username"`
	Password string `json:"password"`
}

func (dto *PlayerConnectionRequest) Validate() error {
	if strings.TrimSpace(dto.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if strings.TrimSpace(dto.Password) == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

type PlayerConnectionResponse struct {
	Token string `json:"token"`
	UserId    string `json:"userId"`
	Username  string `json:"username"`
	Password string `json:"password"`
}

