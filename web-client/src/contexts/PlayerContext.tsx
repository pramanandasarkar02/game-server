import React, { createContext, useState, type ReactNode } from 'react';
import type { Player } from '../types/player';
import { v4 as uuidv4 } from 'uuid';

interface PlayerContextType {
  player: Player | null;
  setPlayer: React.Dispatch<React.SetStateAction<Player | null>>;
}

interface TokenContextType {
  token: string | null;
  setToken: React.Dispatch<React.SetStateAction<string | null>>;
}

export const PlayerContext = createContext<PlayerContextType>({
  player: null,
  setPlayer: () => {},
});

export const TokenContext = createContext<TokenContextType>({
  token: null,
  setToken: () => {},
});

interface AppContextProviderProps {
  children: ReactNode;
}

export const PlayerContextProvider: React.FC<AppContextProviderProps> = ({ children }) => {
  const [player, setPlayer] = useState<Player | null>({
    id: uuidv4(),
    name: 'Guest',
    level: 0,
  });
  const [token, setToken] = useState<string | null>(null);

  return (
    <PlayerContext.Provider value={{ player, setPlayer }}>
      <TokenContext.Provider value={{ token, setToken }}>
        {children}
      </TokenContext.Provider>
    </PlayerContext.Provider>
  );
};