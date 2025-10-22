import { useContext, useState } from "react";
import PlayerContext from "../context/PlayerContext";
import { useNavigate } from "react-router-dom";

const Home = () => {
  const { player } = useContext(PlayerContext);
  const navigate = useNavigate();
  const [selectedGame, setSelectedGame] = useState("snake");

  const findMatchButton = () => {
    navigate("/match-make", { state: { game: selectedGame } });
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

      <button
        onClick={findMatchButton}
        className="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-3 px-8 rounded-lg shadow-md transition-colors"
      >
        Find Match
      </button>
    </div>
  );
};

export default Home;
