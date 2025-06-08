package matchmaking

import (
	"fmt"
	"time"

	"github.com/pramanandasarkar02/game-server/internal/models"
	"github.com/pramanandasarkar02/game-server/internal/store"
	"github.com/pramanandasarkar02/game-server/pkg/logger"
)

func StartMatchmaker(gameStore *store.GameStore, queueStore *store.QueueStore, matchStore *store.MatchStore) {
	for {
		games := gameStore.GetAll()
		for _, g := range games {
			queuedPlayers := queueStore.GetQueue(g.ID())
			if len(queuedPlayers) >= g.RequiredPlayers() {
				matchID := fmt.Sprintf("match-%d", time.Now().UnixNano())
				newMatch := models.Match{
					ID:      matchID,
					GameID:  g.ID(),
					Players: queuedPlayers[:g.RequiredPlayers()],
				}
				queueStore.RemovePlayers(g.ID(), g.RequiredPlayers())
				matchStore.AddMatch(newMatch)
				g.InitializeState(matchID, newMatch.Players)
				logger.Info("Created match %s for game %s with players: %v", matchID, g.ID(), newMatch.Players)
			}
		}
		time.Sleep(1 * time.Second)
	}
}