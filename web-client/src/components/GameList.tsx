import games from '../data/games.json'

interface GameList {
    id: number;
    name: string;
    description: string;
    image: string;
}


type Props = {}

const GameList = (props: Props) => {
  return (
    <div>
        GameList
        <div>
            {games.map((game) => (
                <div key={game.id}>
                    <h2>{game.name}</h2>
                    <p>{game.description}</p>
                    <img src={game.image} alt={game.name} />
                </div>
            ))} 
        </div>

    </div>
  )
}

export default GameList