import { Routes, Route } from 'react-router-dom';
import Home from './pages/Home';
import Games from './pages/Games';

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/games" element={<Games />} />
    </Routes>
  );
}

export default App;