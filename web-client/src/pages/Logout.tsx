import { useContext } from 'react'
import PlayerContext from '../context/PlayerContext'



const Logout = () => {
    const {setPlayer} = useContext(PlayerContext);
    const logoutHandler = () => {
        setPlayer(null)
    }
  return (
    <button onClick={logoutHandler}>logout</button>
  )
}

export default Logout