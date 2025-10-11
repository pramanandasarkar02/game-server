package service

import "log"



type MatchMakeService struct{

}


type PlayerMatchResponse struct{
	PlayerId string `json:"playerId"`
	GameId string	`json:"gameId"`
	MatchId string `json:"matchId"`
	Players []Player	`json:"players"`

}


func NewMatchMakeService() *MatchMakeService{
	return &MatchMakeService{

	}
}


func(ms *MatchMakeService)AddQueue(playerId string, gameId string) error{
	log.Printf("%v request for game %v", playerId, gameId)
	return nil
}

func(ms *MatchMakeService)RemoveQueue(playerId string) error{
	log.Printf("%v removed from the queue", playerId)
	return nil
}


func (ms *MatchMakeService) GetMatch(playerId string) (*PlayerMatchResponse, error) {
	newPlayerResponse := &PlayerMatchResponse{
		PlayerId: playerId,
		GameId:   "abc",
		MatchId:  "123",
		Players:  make([]Player, 0),
	}

	log.Println(newPlayerResponse)
	return newPlayerResponse, nil
}