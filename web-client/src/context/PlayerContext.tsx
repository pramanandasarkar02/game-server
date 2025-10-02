import { createContext, useState } from "react";
import type { ReactNode } from "react";
import type { Player } from "../types/player";

type PlayerContextType = {
  player: Player | null;
  setPlayer: (player: Player | null) => void;
};

const defaultPlayerContext: PlayerContextType = {
  player: {
    username: "",
    userId: "",
    playerStatus: "",
  },

  setPlayer: () => {},
};

const PlayerContext = createContext<PlayerContextType>(defaultPlayerContext);

export const PlayerContextProvider = ({
  children,
}: {
  children: ReactNode;
}) => {
  const [player, setPlayer] = useState<Player | null>(null);

  return (
    <PlayerContext.Provider value={{ player, setPlayer }}>
      {children}
    </PlayerContext.Provider>
  );
};

export default PlayerContext;
