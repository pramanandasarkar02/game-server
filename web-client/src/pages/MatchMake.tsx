import { useNavigate } from "react-router-dom"


const MatchMake = () => {
  const navigate = useNavigate();
  const userId = "12345";
  const gameId = "abcde";
  const goTheGameCanvas = () => {
    navigate(`/snake-game/game-canvas/${gameId}/${userId}`)
  }

    return (

    <div>
        <h1>Snake Game</h1>
        <button onClick={goTheGameCanvas}>
            Enter game
        </button>
        
    </div>
  )
}

export default MatchMake