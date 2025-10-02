import { createContext, useState, useEffect } from "react";
import type { ReactNode } from "react";
import type { Player } from "../types/player";

type PlayerContextType = {
  player: Player | null;
  setPlayer: (player: Player | null) => void;
};

const defaultPlayerContext: PlayerContextType = {
  player: null,
  setPlayer: () => {},
};

const PlayerContext = createContext<PlayerContextType>(defaultPlayerContext);

export const PlayerContextProvider = ({ children }: { children: ReactNode }) => {
  const [player, setPlayerState] = useState<Player | null>(null);

  useEffect(() => {
    const storedPlayer = localStorage.getItem("player");
    if (storedPlayer) {
      setPlayerState(JSON.parse(storedPlayer));
    }
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
    <PlayerContext.Provider value={{ player, setPlayer }}>
      {children}
    </PlayerContext.Provider>
  );
};

export default PlayerContext;
