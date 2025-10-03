import { useContext } from "react"
import PlayerContext from "../context/PlayerContext"
import { useNavigate } from "react-router-dom";

const Home = () => {
    const {player} = useContext(PlayerContext);
    const navigate =  useNavigate();
    const findMatchButton = async () => {
        navigate("/match-make")
    }


  return (
    
    <>
    <h1>Game-server Application</h1>
    <div>
        {player?.username && (
        <p>
          Logged in as <strong>{player.username}</strong> (ID: {player.userId}) && playerStatus: ({player.playerStatus})
        </p>
      )}
    </div>
    <div>
        {/* game List */}
    </div>
    <button onClick={findMatchButton}>
        Find Match
    </button>
    </>
  )
}

export default Home