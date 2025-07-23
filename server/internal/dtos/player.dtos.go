package dtos

import (
	"fmt"
	"strings"
)


type PlayerRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (dto *PlayerRegisterRequest) Validate() error {
	if strings.TrimSpace(dto.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if strings.TrimSpace(dto.Password) == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

type PlayerRegisterResponse struct {
	UserId string `json:"user_id"`
	Username string `json:"username"`
	Token string `json:"token"`	
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
	UserId string `json:"user_id"`
	Username string `json:"username"`
	Token string `json:"token"`
}

type PlayerAuthValidationRequest struct {
	Token string `json:"token"`
}

func (dto *PlayerAuthValidationRequest) Validate() error {
	if strings.TrimSpace(dto.Token) == "" {
		return fmt.Errorf("token is required")
	}
	return nil
}