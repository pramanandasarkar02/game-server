// src/components/Home.tsx
import { useContext, useEffect, useRef, useState } from "react";
import PlayerContext from "../context/PlayerContext";
import { useNavigate } from "react-router-dom";
import axios, { HttpStatusCode } from "axios";
import Logout from "./Logout";
import type { Player } from "../types/player";

type GameEnv = {
  gameId: string;
  matchId: string;
  players: string[];
};

type MatchResponse = {
  isFound: boolean;
  match?: {
    playerId: string;
    gameEnv: GameEnv;
    status: string;
  };
};

const Home = () => {
  const { player, setPlayer } = useContext(PlayerContext);
  const navigate = useNavigate();
  const [selectedGame, setSelectedGame] = useState<string>("snake");
  const [isQueued, setIsQueued] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const pollRef = useRef<number | null>(null);

  const rootUrl = "http://localhost:8080/api";

  /** 🔁 Poll server every 2 seconds to check for a match */
  const startPolling = () => {
    if (pollRef.current) return; // prevent multiple intervals
    const id = window.setInterval(() => {
      getMatch();
    }, 2000);
    pollRef.current = id;
  };

  const stopPolling = () => {
    if (pollRef.current) {
      clearInterval(pollRef.current);
      pollRef.current = null;
    }
  };

  useEffect(() => {
    return () => stopPolling();
  }, []);

  /** 🟢 Add player to queue */
  const addQueue = async () => {
    if (!player?.userId) {
      setError("No player ID found. Please log in.");
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await axios.post(`${rootUrl}/match-make/${player.userId}/${selectedGame}`);

      if (response.status === HttpStatusCode.Ok) {
        console.log("Add queue:", response.data?.message ?? "queued");
        setIsQueued(true);
        startPolling();
      } else {
        setError(`Unexpected response: ${response.status}`);
      }
    } catch (err: any) {
      console.error("Add queue error:", err);
      setError(err?.response?.data?.message ?? err.message ?? "Failed to add to queue");
    } finally {
      setLoading(false);
    }
  };

  /** 🧭 Check if match is found */
  const getMatch = async () => {
    if (!player?.userId) return;

    try {
      const response = await axios.get<MatchResponse>(`${rootUrl}/match-make/${player.userId}`);

      if (response.status === HttpStatusCode.Ok) {
        const data = response.data;
        if (data.isFound && data.match) {
          console.log(`✅ Match found: ${data.match.gameEnv.matchId}`);
          const newPlayer: Player = {
            username: player.username,
            userId: player.userId,
            playerStatus: player.playerStatus,
            matchId: data.match.gameEnv.matchId
          };
          setPlayer(newPlayer);
          // Stop polling once found
          stopPolling();

          // Navigate to the match page with match info
          navigate("/match-make", {
            state: {
              game: data.match.gameEnv.gameId,
              match: data.match.gameEnv,
            },
          });
        } else {
          console.log("⏳ Waiting for match...");
        }
      }
    } catch (err: any) {
      console.error("Error checking match:", err);
    }
  };

  /**  Remove player from queue */
  const removeQueue = async () => {
    if (!player?.userId) return;
    stopPolling();
    setLoading(true);
    try {
      const response = await axios.patch(`${rootUrl}/match-make/${player.userId}`);
      if (response.status === HttpStatusCode.Ok) {
        console.log(response.data?.message ?? "Removed from queue");
        setIsQueued(false);
      }
    } catch (err: any) {
      console.error("Remove queue error:", err);
    } finally {
      setLoading(false);
    }
  };

  /**  Add or remove queue button handler */
  const findMatchButton = () => {
    if (isQueued) removeQueue();
    else addQueue();
  };

  return (
    <div className="min-h-screen bg-gray-900 text-gray-100 flex flex-col items-center justify-center p-6">
      <h1 className="text-4xl font-bold mb-6">Game Server Application</h1>

      {player?.username && (
        <div className="bg-gray-800 shadow-md rounded-lg p-4 mb-6 w-full max-w-md text-center">
          <p className="text-gray-100 text-lg">
            Logged in as <strong>{player.username}</strong> (ID: {player.userId}) <br />
            Status: <span className="text-blue-400">{player.playerStatus}</span>
          </p>
        </div>
      )}

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 mb-6 w-full max-w-2xl">
        {/* Snake Game */}
        <div
          onClick={() => setSelectedGame("snake")}
          className={`cursor-pointer rounded-lg p-6 flex flex-col items-center justify-center shadow-md transition-transform ${
            selectedGame === "snake"
              ? "bg-blue-700 scale-105 border-2 border-blue-400"
              : "bg-gray-800 hover:bg-gray-700"
          }`}
        >
          <h2 className="text-2xl font-semibold mb-2">Snake Game</h2>
          <p className="text-gray-300 text-center">
            Classic snake game. Compete with other players!
          </p>
        </div>

        {/* Tic-Tac-Toe Game */}
        <div
          onClick={() => setSelectedGame("tic-tac-toe")}
          className={`cursor-pointer rounded-lg p-6 flex flex-col items-center justify-center shadow-md transition-transform ${
            selectedGame === "tic-tac-toe"
              ? "bg-blue-700 scale-105 border-2 border-blue-400"
              : "bg-gray-800 hover:bg-gray-700"
          }`}
        >
          <h2 className="text-2xl font-semibold mb-2">Tic-Tac-Toe</h2>
          <p className="text-gray-300 text-center">Play Tic-Tac-Toe with friends!</p>
        </div>
      </div>

      <div className="flex items-center gap-4">
        <button
          onClick={findMatchButton}
          disabled={loading}
          className="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-3 px-8 rounded-lg shadow-md transition-colors disabled:opacity-50"
        >
          {isQueued ? "Cancel Queue" : "Add Queue"}
        </button>

        {isQueued && (
          <button
            onClick={removeQueue}
            disabled={loading}
            className="bg-red-600 hover:bg-red-700 text-white font-semibold py-3 px-6 rounded-lg shadow-md transition-colors disabled:opacity-50"
          >
            Leave Queue
          </button>
        )}
      </div>

      {error && <p className="text-red-400 mt-4">{error}</p>}

      <div className="mt-6">
        <Logout />
      </div>
    </div>
  );
};

export default Home;
