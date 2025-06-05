package models

type Player struct{
	ID string `json:"id"`
	Name string `json:"name"`
	Level float32 `json:"level"`
}