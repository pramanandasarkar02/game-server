import games from '../data/games.json'

interface Game {
    id: number;
    name: string;
    description: string;
    image: string;
}

const GameList = () => {
  return (
    <div className="container mx-auto py-8">
        <h1 className="text-3xl font-bold mb-8 text-center">Available Games</h1>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {games.map((game) => (
                <div 
                    key={game.id} 
                    className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition cursor-pointer"
                >
                    <img 
                        src={game.image} 
                        alt={game.name} 
                        className="w-full h-48 object-cover"
                    />
                    <div className="p-4">
                        <h2 className="text-xl font-semibold mb-2">{game.name}</h2>
                        <p className="text-gray-600">{game.description}</p>
                        <button className="mt-4 bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition">
                            Play Now
                        </button>
                    </div>
                </div>
            ))} 
        </div>
    </div>
  )
}

export default GameList