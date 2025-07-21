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
	Success bool   `json:"success"`
	Message string `json:"message"`
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