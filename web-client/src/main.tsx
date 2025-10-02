import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import { PlayerContextProvider } from './context/PlayerContext.tsx'

createRoot(document.getElementById('root')!).render(
    <PlayerContextProvider>
      <App />
    </PlayerContextProvider>
)
