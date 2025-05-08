import React, { createContext, useContext, useState } from 'react';
import { Game, Lobby, LobbyPlayer } from '../types';
import { useAuth } from './AuthContext';

interface GamesContextType {
  games: Game[];
  lobbies: Lobby[];
  currentLobby: Lobby | null;
  joinLobby: (gameId: string) => void;
  leaveLobby: () => void;
  setPlayerReady: (isReady: boolean) => void;
}

// Mock data
const MOCK_GAMES: Game[] = [
  {
    id: 'tictactoe',
    name: 'Tic Tac Toe',
    description: 'Classic game of X and O. Be the first to get three in a row!',
    image: 'https://images.pexels.com/photos/5887520/pexels-photo-5887520.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=1',
    minPlayers: 2,
    maxPlayers: 2,
    difficulty: 'easy'
  },
  {
    id: 'puzzle',
    name: 'Sliding Puzzle',
    description: 'Arrange the tiles in numerical order by sliding them into the empty space.',
    image: 'https://images.pexels.com/photos/3165335/pexels-photo-3165335.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=1',
    minPlayers: 1,
    maxPlayers: 1,
    difficulty: 'medium'
  },
  {
    id: 'memory',
    name: 'Memory Match',
    description: 'Test your memory by matching pairs of cards.',
    image: 'https://images.pexels.com/photos/3368816/pexels-photo-3368816.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=1',
    minPlayers: 1,
    maxPlayers: 2,
    difficulty: 'medium'
  },
  {
    id: 'chess',
    name: 'Chess',
    description: 'The classic strategy board game of kings and queens.',
    image: 'https://images.pexels.com/photos/260024/pexels-photo-260024.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=1',
    minPlayers: 2,
    maxPlayers: 2,
    difficulty: 'hard'
  }
];

const MOCK_LOBBIES: Lobby[] = [
  {
    id: 'lobby1',
    gameId: 'tictactoe',
    gameName: 'Tic Tac Toe',
    players: [
      { id: '1', username: 'player1', level: 5, isReady: true },
    ],
    maxPlayers: 2,
    status: 'waiting',
    createdAt: new Date()
  },
  {
    id: 'lobby2',
    gameId: 'puzzle',
    gameName: 'Sliding Puzzle',
    players: [
      { id: '2', username: 'player2', level: 3, isReady: false },
    ],
    maxPlayers: 1,
    status: 'waiting',
    createdAt: new Date()
  }
];

const GamesContext = createContext<GamesContextType | undefined>(undefined);

export const GamesProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user } = useAuth();
  const [games] = useState<Game[]>(MOCK_GAMES);
  const [lobbies, setLobbies] = useState<Lobby[]>(MOCK_LOBBIES);
  const [currentLobby, setCurrentLobby] = useState<Lobby | null>(null);

  const joinLobby = (gameId: string) => {
    if (!user) return;

    // Find existing lobby for this game with space
    let lobby = lobbies.find(l => 
      l.gameId === gameId && 
      l.players.length < l.maxPlayers &&
      l.status === 'waiting'
    );

    // If no lobby exists, create a new one
    if (!lobby) {
      const game = games.find(g => g.id === gameId);
      if (!game) return;

      const newLobby: Lobby = {
        id: `lobby-${Date.now()}`,
        gameId,
        gameName: game.name,
        players: [],
        maxPlayers: game.maxPlayers,
        status: 'waiting',
        createdAt: new Date()
      };
      
      setLobbies(prev => [...prev, newLobby]);
      lobby = newLobby;
    }

    // Add player to lobby
    const playerToAdd: LobbyPlayer = {
      id: user.id,
      username: user.username,
      level: user.level,
      isReady: false
    };

    const updatedLobby = {
      ...lobby,
      players: [...lobby.players, playerToAdd]
    };

    setLobbies(prev => 
      prev.map(l => l.id === updatedLobby.id ? updatedLobby : l)
    );
    setCurrentLobby(updatedLobby);
  };

  const leaveLobby = () => {
    if (!currentLobby || !user) return;

    // Remove player from lobby
    const updatedLobby = {
      ...currentLobby,
      players: currentLobby.players.filter(p => p.id !== user.id)
    };

    // If lobby is empty, remove it
    if (updatedLobby.players.length === 0) {
      setLobbies(prev => prev.filter(l => l.id !== currentLobby.id));
    } else {
      setLobbies(prev => 
        prev.map(l => l.id === updatedLobby.id ? updatedLobby : l)
      );
    }

    setCurrentLobby(null);
  };

  const setPlayerReady = (isReady: boolean) => {
    if (!currentLobby || !user) return;

    // Update player ready status
    const updatedPlayers = currentLobby.players.map(player => 
      player.id === user.id ? { ...player, isReady } : player
    );

    const updatedLobby = {
      ...currentLobby,
      players: updatedPlayers,
      // If all players are ready and we have minimum players, set status to starting
      status: updatedPlayers.every(p => p.isReady) && 
              updatedPlayers.length >= games.find(g => g.id === currentLobby.gameId)?.minPlayers!
              ? 'starting' : 'waiting'
    };

    setLobbies(prev => 
      prev.map(l => l.id === updatedLobby.id ? updatedLobby : l)
    );
    setCurrentLobby(updatedLobby);
  };

  return (
    <GamesContext.Provider
      value={{
        games,
        lobbies,
        currentLobby,
        joinLobby,
        leaveLobby,
        setPlayerReady
      }}
    >
      {children}
    </GamesContext.Provider>
  );
};

export const useGames = (): GamesContextType => {
  const context = useContext(GamesContext);
  if (context === undefined) {
    throw new Error('useGames must be used within a GamesProvider');
  }
  return context;
};