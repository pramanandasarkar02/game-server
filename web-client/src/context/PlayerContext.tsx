// context/PlayerContext.tsx
import { createContext, useState, useEffect } from "react";
import type { ReactNode } from "react";
import type { Player } from "../types/player";

type PlayerContextType = {
  player: Player | null;
  setPlayer: (player: Player | null) => void;
  loading: boolean; 
};

const defaultPlayerContext: PlayerContextType = {
  player: null,
  setPlayer: () => {},
  loading: true,
};

const PlayerContext = createContext<PlayerContextType>(defaultPlayerContext);

export const PlayerContextProvider = ({ children }: { children: ReactNode }) => {
  const [player, setPlayerState] = useState<Player | null>(null);
  const [loading, setLoading] = useState(true); 

  useEffect(() => {
    const storedPlayer = localStorage.getItem("player");
    if (storedPlayer) {
      try {
        setPlayerState(JSON.parse(storedPlayer));
      } catch (error) {
        console.error("Failed to parse player from localStorage", error);
        localStorage.removeItem("player");
      }
    }
    setLoading(false); 
  }, []);

  const setPlayer = (player: Player | null) => {
    if (player) {
      localStorage.setItem("player", JSON.stringify(player));
    } else {
      localStorage.removeItem("player");
    }
    setPlayerState(player);
  };

  return (
    <PlayerContext.Provider value={{ player, setPlayer, loading }}>
      {children}
    </PlayerContext.Provider>
  );
};

export default PlayerContext;