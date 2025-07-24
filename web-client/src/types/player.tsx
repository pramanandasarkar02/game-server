export interface Player {
  id: string;
  name: string;
  level: number;
}

export interface ChatMessage {
  matchId: string;
  playerId: string;
  content: string;
  timestamp: string;
}

export interface TicTacToeState {
  board: string[];
  turn: string;
  winner: string;
  isDraw: boolean;
  players: string[];
}

export interface WebSocketMessage {
  type: 'chat' | 'move' | 'state';
  data: ChatMessage | { index: number; playerID: string } | TicTacToeState;
}


export interface PlayerInformation {
  id: string;
  userName: string;
  level: number;
  matchHistory: string[];
  score: number;
}