package handler

import (
	"fmt"
	"net/http"
)

type PlayerHandler struct{}

func NewPlayerHandler() *PlayerHandler {
	return &PlayerHandler{}
}

func (ph *PlayerHandler) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Player Login called")
}

func (ph *PlayerHandler) Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Player Logout called")
}

func (ph *PlayerHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Player Signup called")
}
