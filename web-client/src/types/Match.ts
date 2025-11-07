export type GameEnv = {
    gameId: string,
    matchId: string,
    players: string[],
}

export type MatchContext = {
    playerId: string,
    gameEnv: GameEnv,
    status: string
}