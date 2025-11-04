// src/components/Home.tsx
import { useContext, useEffect, useRef, useState } from "react";
import PlayerContext from "../context/PlayerContext";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import Logout from "./Logout";

type GameEnv = {
  gameId: string;
  matchId: string;
  players: string[];
};

type MatchCheckResponse = {
  isFound?: boolean;
  gameEnv?: GameEnv;
  message?: string;
};

const Home = () => {
  const { player } = useContext(PlayerContext);
  const navigate = useNavigate();
  const [selectedGame, setSelectedGame] = useState<string>("snake");
  const [isQueued, setIsQueued] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const pollRef = useRef<number | null>(null);

  const rootUrl = "http://localhost:8080/api";

  // Start polling for a match (every 2s) and keep the interval id in pollRef.
  const startPolling = () => {
    // avoid double intervals
    if (pollRef.current) return;
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
    // cleanup on unmount
    return () => {
      stopPolling();
    };
  }, []);

  const addQueue = async () => {
    if (!player?.userId) {
      setError("No player ID found. Please login.");
      return;
    }

    setLoading(true);
    setError(null);

    const requestString = `${rootUrl}/match-make/${player.userId}/${selectedGame}`;

    try {
      const response = await axios.post(requestString);

      // Expect backend to return 200 + { message }
      if (response.status === 200) {
        const data = response.data;
        console.log("Add queue:", data?.message ?? "queued");
        setIsQueued(true);
        // start polling for match
        startPolling();
      } else {
        setError(`Unexpected response status: ${response.status}`);
      }
    } catch (err: any) {
      console.error("Failed to add to queue:", err);
      setError(err?.response?.data?.message ?? err.message ?? "Failed to add to queue");
    } finally {
      setLoading(false);
    }
  };

  const getMatch = async () => {
    // Polling - check if player has been matched
    if (!player?.userId) {
      stopPolling();
      setIsQueued(false);
      return;
    }

    const requestString = `${rootUrl}/match-make/${player.userId}`;

    try {
      const response = await axios.patch(requestString);
      if (response.status === 200) {
        const data: MatchCheckResponse = response.data ?? {};

        // If backend returns an explicit isFound boolean and a gameEnv, navigate
        if (data.isFound && data.gameEnv) {
          console.log(`Match found: ${data.gameEnv.matchId}`, data.gameEnv);
          stopPolling();
          setIsQueued(false);

          // pass full gameEnv to match-make route for the next screen
          navigate("/match-make", { state: { gameEnv: data.gameEnv } });
        } else {
          // still queued - optional console
          console.log("No match yet.");
        }
      } else {
        console.warn("Unexpected getMatch status:", response.status);
      }
    } catch (err: any) {
      // If server replies with no active match or other error, just log and keep polling.
      console.error("Error checking match:", err?.response?.data ?? err.message ?? err);
      // If the backend says player not found or similar, you may want to stop polling:
      // Example: if err.response?.status === 404 -> stopPolling()
    }
  };

  const removeQueue = async () => {
    if (!player?.userId) {
      setError("No player ID found. Please login.");
      return;
    }

    setLoading(true);
    setError(null);

    const requestString = `${rootUrl}/match-make/${player.userId}`;

    try {
      // You used PATCH originally — keep that so we don't change backend expectations.
      const response = await axios.patch(requestString);
      if (response.status === 200) {
        const data = response.data;
        console.log("Removed from queue:", data?.message ?? "removed");
        setIsQueued(false);
        stopPolling();
      } else {
        setError(`Unexpected response status: ${response.status}`);
      }
    } catch (err: any) {
      console.error("Failed to remove from queue:", err);
      setError(err?.response?.data?.message ?? err.message ?? "Failed to remove from queue");
    } finally {
      setLoading(false);
    }
  };

  // click handler for Add queue button
  const findMatchButton = () => {
    // toggle behavior: if already queued, remove; else add
    if (isQueued) {
      removeQueue();
    } else {
      addQueue();
    }
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
          <p className="text-gray-300 text-center">Classic snake game. Compete with other players!</p>
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
          {isQueued ? "Cancel Queue" : "Add queue"}
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
