import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import Login from './pages/Login'
import { useContext } from 'react'
import PlayerContext from './context/PlayerContext'
import Home from './pages/Home'
import Signup from './pages/Signup'

function App() {
  const {player} = useContext(PlayerContext);
  return (
    <BrowserRouter>
    <Routes>
      <Route path='/' element = {player ? <Home /> :<Navigate to="/login" />}/>
      <Route path='/login' element ={<Login />}/>
      <Route path='/signup' element ={<Signup />}/>
    </Routes>
    </BrowserRouter>
  )
}

export default App
