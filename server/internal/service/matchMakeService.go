package service

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)


var(
	MINIMUM_MATCH_PLAYER = 2
)


type MatchMakeService struct{
	queue map[string]string 	//playerId -> gameId
	matches map[string]GameEnv  //playerId -> gameEnv
}

type GameEnv struct{
	GameId string `json:"gameId"`
	MatchId string `json:"matchId"`
	Players []string `json:"players"`
}


type PlayerMatchResponse struct{
	PlayerId string `json:"playerId"`
	GameEnv GameEnv `json:"gameEnv"`

}


func NewMatchMakeService() *MatchMakeService{
	return &MatchMakeService{
		queue: make(map[string]string),
		matches: make(map[string]GameEnv),
	}
}


func(ms *MatchMakeService)AddQueue(playerId string, gameId string) error{
	// log.Printf("%v request for game %v", playerId, gameId)
	// return nil

	if _, exists := ms.queue[playerId]; exists{
		return fmt.Errorf("player already in queue with game id %v", ms.queue[playerId])
	}

	ms.queue[playerId] = gameId
	log.Printf("Player %v added to queue for game %v", playerId, gameId)

	ms.matchMake(gameId)
	return nil	
}

func(ms *MatchMakeService)RemoveQueue(playerId string) error{
	if _, exists := ms.queue[playerId]; !exists{
		return fmt.Errorf("player %v not found in queue", playerId)
	}
	delete(ms.queue, playerId)
	
	
	log.Printf("%v removed from the queue", playerId)
	return nil
}


func (ms *MatchMakeService) GetMatch(playerId string) (*PlayerMatchResponse, error) {

	if _, exists := ms.matches[playerId]; !exists{
		return &PlayerMatchResponse{}, fmt.Errorf("playerId %v not found in the current matches", playerId)
	}

	gameEnv := ms.matches[playerId]

	
	newPlayerResponse := &PlayerMatchResponse{
		PlayerId: playerId,
		GameEnv: gameEnv,
		
	}

	log.Println(newPlayerResponse)
	return newPlayerResponse, nil
}


func (ms *MatchMakeService) matchMake(gameId string){
	players := []string{}

	for playerId, gId := range ms.queue {
		if gId == gameId{
			players = append(players, playerId)
		}
	}
	if len(players) >= MINIMUM_MATCH_PLAYER{
		matchId := fmt.Sprintf("match-%v", uuid.New())
		log.Printf("Creating match %v for game %v with players: %v", matchId, gameId, players)
		
		gameEnv := GameEnv{
			GameId: gameId,
			MatchId: matchId,
			Players: players, 
		}

		for _, p := range players{
			delete(ms.queue, p)
			ms.matches[p] = gameEnv
		}
		log.Printf("Match %v created successfully", matchId)
	}
}

