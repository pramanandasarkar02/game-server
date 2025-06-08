import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './App';
import './index.css';
import { PlayerContextProvider } from './contexts/PlayerContext';

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <PlayerContextProvider>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </PlayerContextProvider>
  </StrictMode>,
);