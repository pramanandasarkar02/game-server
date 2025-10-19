import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import Login from './pages/Login'
import { useContext } from 'react'
import PlayerContext from './context/PlayerContext'
import Home from './pages/Home'
import Signup from './pages/Signup'
// import SnakeGame from './pages/SnakeGame'
import MatchMake from './pages/MatchMake'
import SnakeGame from './pages/snakeGame/SnakeGame'

function App() {
  const {player} = useContext(PlayerContext);
  return (
    <BrowserRouter>
    <Routes>
      <Route path='/' element = {player ? <Home /> :<Navigate to="/login" />}/>
      <Route path='/login' element ={<Login />}/>
      <Route path='/signup' element ={<Signup />}/>
      {/* <Route path='/snake-game/game-canvas/:gameId/:playerId' element={<SnakeGame />} /> */}
      <Route path='/match-make' element={<MatchMake />} />
      <Route path='/game' element={<SnakeGame />} />
    </Routes>
    </BrowserRouter>
  )
}

export default App
