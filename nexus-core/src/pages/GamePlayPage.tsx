import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeft, Trophy, Users } from 'lucide-react';
import Layout from '../components/layout/Layout';
import Button from '../components/ui/Button';
import Card from '../components/ui/Card';
import TicTacToeBoard from '../games/tic-tac-toe/TicTacToeBoard';
import PuzzleBoard from '../games/puzzle/PuzzleBoard';
import { useGames } from '../contexts/GamesContext';
import { useAuth } from '../contexts/AuthContext';

const GamePlayPage: React.FC = () => {
  const { gameId } = useParams<{ gameId: string }>();
  const { games } = useGames();
  const { user } = useAuth();
  const navigate = useNavigate();
  const [game, setGame] = useState(games.find(g => g.id === gameId));
  const [gameResult, setGameResult] = useState<'win' | 'lose' | 'draw' | null>(null);
  const [earnedPoints, setEarnedPoints] = useState(0);

  useEffect(() => {
    if (!gameId || !game) {
      navigate('/games');
    }
  }, [gameId, game, navigate]);

  const handleGameEnd = (result: 'win' | 'lose' | 'draw' | number) => {
    // Handle number result from puzzle game (moves)
    if (typeof result === 'number') {
      const points = Math.max(100, 500 - (result * 5));
      setEarnedPoints(points);
      setGameResult('win');
      return;
    }
    
    // Handle win/lose/draw result
    setGameResult(result);
    
    // Calculate points earned
    let points = 0;
    if (result === 'win') {
      points = 100;
    } else if (result === 'draw') {
      points = 25;
    }
    
    setEarnedPoints(points);
  };

  const handlePlayAgain = () => {
    setGameResult(null);
    setEarnedPoints(0);
  };

  const handleReturnToLobby = () => {
    navigate('/games');
  };

  if (!game) return null;

  const renderGame = () => {
    switch (gameId) {
      case 'tictactoe':
        return <TicTacToeBoard onGameEnd={handleGameEnd} />;
      case 'puzzle':
        return <PuzzleBoard onGameEnd={handleGameEnd} />;
      default:
        return (
          <div className="text-center py-16">
            <h3 className="text-xl">Game not available</h3>
            <p className="text-gray-400 mt-2">This game is currently under development</p>
          </div>
        );
    }
  };

  return (
    <Layout>
      <div className="max-w-4xl mx-auto">
        <div className="flex justify-between items-center mb-6">
          <div className="flex items-center">
            <button
              onClick={handleReturnToLobby}
              className="p-2 mr-2 rounded-full hover:bg-gray-800 transition-colors"
            >
              <ArrowLeft className="h-5 w-5" />
            </button>
            <h1 className="text-2xl font-bold">{game.name}</h1>
          </div>
          
          <div className="flex items-center space-x-4">
            <div className="flex items-center text-gray-400">
              <Users className="h-4 w-4 mr-1" />
              <span className="text-sm">{game.minPlayers}-{game.maxPlayers} players</span>
            </div>
            
            <div className="hidden md:flex items-center px-3 py-1 rounded-full bg-primary-900/20 text-primary-300">
              <Trophy className="h-4 w-4 mr-1" />
              <span className="text-sm">{user?.level || 0}</span>
            </div>
          </div>
        </div>
        
        {/* Game result modal */}
        {gameResult && (
          <div className="fixed inset-0 flex items-center justify-center z-50 bg-black/70">
            <Card className="w-full max-w-md p-6 animate-float">
              <div className="text-center mb-6">
                <div className={`
                  w-20 h-20 mx-auto rounded-full flex items-center justify-center mb-4
                  ${gameResult === 'win' ? 'bg-green-900/30' : 
                    gameResult === 'draw' ? 'bg-yellow-900/30' : 'bg-red-900/30'}
                `}>
                  <Trophy className={`h-10 w-10 
                    ${gameResult === 'win' ? 'text-green-400' : 
                      gameResult === 'draw' ? 'text-yellow-400' : 'text-red-400'}
                  `} />
                </div>
                
                <h2 className="text-2xl font-bold mb-2">
                  {gameResult === 'win' ? 'Victory!' : 
                   gameResult === 'draw' ? 'Draw!' : 'Defeat!'}
                </h2>
                
                <p className="text-gray-300 mb-4">
                  {gameResult === 'win' ? 'Congratulations on your win!' : 
                   gameResult === 'draw' ? 'So close! The game ended in a draw.' : 
                   'Better luck next time!'}
                </p>
                
                {earnedPoints > 0 && (
                  <div className="bg-primary-900/20 rounded-lg p-3 mb-6">
                    <p className="text-primary-300 font-semibold">
                      + {earnedPoints} points earned!
                    </p>
                  </div>
                )}
                
                <div className="flex flex-col sm:flex-row gap-3">
                  <Button
                    onClick={handlePlayAgain}
                    fullWidth
                  >
                    Play Again
                  </Button>
                  <Button
                    variant="outline"
                    onClick={handleReturnToLobby}
                    fullWidth
                  >
                    Return to Lobby
                  </Button>
                </div>
              </div>
            </Card>
          </div>
        )}
        
        <div className="bg-background-card rounded-lg p-6">
          {renderGame()}
        </div>
      </div>
    </Layout>
  );
};

export default GamePlayPage;