import { Routes, Route, useLocation } from 'react-router-dom';
import Home from './pages/Home';
import Games from './pages/Games';
import PlayerHome from './pages/PlayerHome';
import Profile from './pages/Profile';
import Friends from './pages/Friends';
import Header from './components/Header';
import Footer from './components/Footer';

function App() {
  const location = useLocation();
  const isHomePage = location.pathname === '/';

  return (
    <div className="flex flex-col min-h-screen">
      {!isHomePage && <Header />}
      <main className="flex-grow">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/home" element={<PlayerHome />} />
          <Route path="/games" element={<Games />} />
          <Route path="/profile" element={<Profile />} />
          <Route path="/friends" element={<Friends />} />
        </Routes>
      </main>
      <Footer />
    </div>
  );
}

export default App;