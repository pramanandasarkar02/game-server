import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { useContext } from "react";
import PlayerContext from "./context/PlayerContext";

import Login from "./pages/Login";
import Signup from "./pages/Signup";
import Home from "./pages/Home";
import MatchMake from "./pages/MatchMake";
import SnakeGame from "./pages/snakeGame/SnakeGame";

function App() {
  const { player, loading } = useContext(PlayerContext);

  // Show nothing or a loader while checking auth
  if (loading) {
    return <div>Loading...</div>; // Or a spinner
  }

  return (
    <BrowserRouter>
      <Routes>
        {/* Home: protected */}
        <Route path="/" element={player ? <Home /> : <Navigate to="/login" />} />

        {/* Auth routes */}
        <Route path="/login" element={!player ? <Login /> : <Navigate to="/" />} />
        <Route path="/signup" element={!player ? <Signup /> : <Navigate to="/" />} />

        {/* Matchmaking: protected */}
        <Route
          path="/match-make"
          element={player ? <MatchMake /> : <Navigate to="/login" />}
        />

        {/* Snake game */}
        <Route
          path="/snake-game/game-canvas/:gameId/:playerId"
          element={player ? <SnakeGame /> : <Navigate to="/login" />}
        />

        {/* Fallback */}
        <Route path="*" element={<Navigate to={player ? "/" : "/login"} />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;