package shared

import "encoding/json"


type Message struct {
	Type string 	`json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type User struct{
	ID string `json:"id"`
	Game string `json:"username"`
	Token string `json:"-"`
}

type GameRoom struct{
	ID string `json:"id"`
	Game string `json:"game"`
	Players []string `json:"player"`
	Status string `json:"status"` // "waiting", "playing", "finished"
}


const(
	MsgTypeAuth = "auth"
	MsgTypeJoinRoom = "join_room"
	MsgTypeLeaveRoom = "leave_room"
	MsgTypeGameMove = "game_move"
	MsgTypeChat = "chat"
)


