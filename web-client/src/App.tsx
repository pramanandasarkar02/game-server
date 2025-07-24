import { Routes, Route } from 'react-router-dom';
import Home from './pages/Home';
import Games from './pages/Games';
import PlayerHome from './pages/PlayerHome';
import Profile from './pages/Profile';

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/home" element={<PlayerHome />} />
      <Route path="/games" element={<Games />} />
      <Route path="/profile" element={<Profile />} />

    </Routes>
  );
}

export default App;