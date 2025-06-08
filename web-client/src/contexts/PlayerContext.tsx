import React, { createContext, useState, type ReactNode } from 'react';
import type { Player } from '../types/player';
import { v4 as uuidv4 } from 'uuid';

interface PlayerContextType {
  player: Player | null;
  setPlayer: React.Dispatch<React.SetStateAction<Player | null>>;
}

export const PlayerContext = createContext<PlayerContextType>({
  player: null,
  setPlayer: () => {},
});

interface PlayerContextProviderProps {
  children: ReactNode;
}

export const PlayerContextProvider: React.FC<PlayerContextProviderProps> = ({ children }) => {
  const [player, setPlayer] = useState<Player | null>({
    id: uuidv4(),
    name: 'Guest',
    level: 0,
  });

  return (
    <PlayerContext.Provider value={{ player, setPlayer }}>
      {children}
    </PlayerContext.Provider>
  );
};