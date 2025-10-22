import { useNavigate, useLocation } from "react-router-dom";

const MatchMake = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const game = (location.state as { game: string })?.game || "snake";
  const userId = "12345";
  const gameId = "abcde";

  const goToGameCanvas = () => {
    if (game === "snake") {
      navigate(`/snake-game/game-canvas/${gameId}/${userId}`);
    } else if (game === "tic-tac-toe") {
      navigate(`/tic-tac-toe/${gameId}/${userId}`);
    }
  };

  return (
    <div className="min-h-screen bg-gray-900 flex flex-col items-center justify-center p-6 text-gray-100">
      <h1 className="text-4xl font-bold mb-6">{game === "snake" ? "Snake Game" : "Tic-Tac-Toe"}</h1>
      <div className="bg-gray-800 p-8 rounded-xl shadow-lg w-full max-w-md text-center">
        <p className="text-gray-300 mb-6">
          Ready to play {game}? Click below to enter the game.
        </p>
        <button
          onClick={goToGameCanvas}
          className="bg-green-600 hover:bg-green-700 text-white font-semibold py-3 px-8 rounded-lg shadow-md transition-all transform hover:scale-105"
        >
          Enter Game
        </button>
      </div>
      <p className="mt-6 text-gray-500 text-sm">
        Game ID: <span className="text-blue-400">{gameId}</span> | User ID: <span className="text-blue-400">{userId}</span>
      </p>
    </div>
  );
};

export default MatchMake;
