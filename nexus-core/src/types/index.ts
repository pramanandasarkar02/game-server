export interface User {
  id: string;
  username: string;
  email?: string;
  level: number;
  score: number;
  gamesPlayed: number;
  gamesWon: number;
}

export interface Game {
  id: string;
  name: string;
  description: string;
  image: string;
  minPlayers: number;
  maxPlayers: number;
  difficulty: 'easy' | 'medium' | 'hard';
}

export interface LobbyPlayer {
  id: string;
  username: string;
  level: number;
  isReady: boolean;
}

export interface Lobby {
  id: string;
  gameId: string;
  gameName: string;
  players: LobbyPlayer[];
  maxPlayers: number;
  status: 'waiting' | 'starting' | 'playing' | 'finished';
  createdAt: Date;
}

export interface TicTacToeGame {
  board: Array<string | null>;
  currentPlayer: string;
  winner: string | null;
  isDraw: boolean;
}

export interface PuzzleGame {
  tiles: number[];
  emptyIndex: number;
  size: number;
  moves: number;
  isComplete: boolean;
}