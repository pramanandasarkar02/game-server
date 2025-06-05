package models

type Game struct{
	ID string `json:"id"`
	Title string `json:"title"`
	RequiredPlayer int `json:"requiredPlayer"`
}