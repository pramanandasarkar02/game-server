package models

import "fmt"

type Player struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Level float32 `json:"level"`
}

func (p *Player) Validate() error {
	if p.ID == "" || p.Name == "" {
		return fmt.Errorf("player ID and name are required")
	}
	return nil
}