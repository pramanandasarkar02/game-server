import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './App';
import './index.css';
import { PlayerContextProvider } from './contexts/PlayerContext';
import { AuthContextProvider } from './contexts/AuthContext';

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AuthContextProvider>
    <PlayerContextProvider>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </PlayerContextProvider>
    </AuthContextProvider>
  </StrictMode>,
);