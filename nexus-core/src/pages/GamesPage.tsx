import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Search } from 'lucide-react';
import Layout from '../components/layout/Layout';
import GameCard from '../components/games/GameCard';
import { useGames } from '../contexts/GamesContext';

const GamesPage: React.FC = () => {
  const { games, joinLobby } = useGames();
  const navigate = useNavigate();
  const [searchTerm, setSearchTerm] = React.useState('');

  const handleGameSelect = (gameId: string) => {
    joinLobby(gameId);
    navigate(`/lobby`);
  };

  const filteredGames = games.filter(game => 
    game.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    game.description.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <Layout>
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">Game Library</h1>
        <p className="text-gray-400">Select a game to join a lobby and start playing!</p>
      </div>

      <div className="relative mb-6">
        <input
          type="text"
          placeholder="Search games..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="w-full bg-background-light border border-gray-700 rounded-lg py-3 px-4 pl-12 text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
        />
        <Search className="absolute left-4 top-3.5 h-5 w-5 text-gray-400" />
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {filteredGames.map(game => (
          <GameCard
            key={game.id}
            game={game}
            onClick={() => handleGameSelect(game.id)}
          />
        ))}
      </div>

      {filteredGames.length === 0 && (
        <div className="text-center py-16">
          <h3 className="text-xl font-medium">No games found</h3>
          <p className="text-gray-400 mt-2">Try a different search term</p>
        </div>
      )}
    </Layout>
  );
};

export default GamesPage;