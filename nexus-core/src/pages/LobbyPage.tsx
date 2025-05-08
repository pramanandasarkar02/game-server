import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Clock, ArrowLeft } from 'lucide-react';
import Layout from '../components/layout/Layout';
import Button from '../components/ui/Button';
import Card from '../components/ui/Card';
import { useGames } from '../contexts/GamesContext';
import { useAuth } from '../contexts/AuthContext';

const LobbyPage: React.FC = () => {
  const { currentLobby, leaveLobby, setPlayerReady } = useGames();
  const { user } = useAuth();
  const navigate = useNavigate();
  const [startCountdown, setStartCountdown] = useState<number | null>(null);
  
  // Redirect if no lobby is selected
  useEffect(() => {
    if (!currentLobby) {
      navigate('/games');
    }
  }, [currentLobby, navigate]);

  // Start countdown if game is starting
  useEffect(() => {
    if (currentLobby?.status === 'starting') {
      setStartCountdown(5);
      
      const timer = setInterval(() => {
        setStartCountdown(prev => {
          if (prev === null || prev <= 1) {
            clearInterval(timer);
            navigate(`/play/${currentLobby.gameId}`);
            return null;
          }
          return prev - 1;
        });
      }, 1000);
      
      return () => clearInterval(timer);
    } else {
      setStartCountdown(null);
    }
  }, [currentLobby?.status, navigate]);

  if (!currentLobby) return null;

  const currentPlayer = currentLobby.players.find(player => player.id === user?.id);
  const isPlayerReady = currentPlayer?.isReady || false;

  return (
    <Layout>
      <div className="max-w-2xl mx-auto">
        <div className="flex items-center mb-6">
          <button
            onClick={() => navigate('/games')}
            className="p-2 mr-2 rounded-full hover:bg-gray-800 transition-colors"
          >
            <ArrowLeft className="h-5 w-5" />
          </button>
          <h1 className="text-2xl font-bold">{currentLobby.gameName} Lobby</h1>
        </div>

        <Card className="mb-6">
          <div className="p-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-semibold">Lobby Status</h2>
              <div className={`
                flex items-center px-3 py-1 rounded-full
                ${currentLobby.status === 'waiting' ? 'bg-yellow-800/30 text-yellow-300' : 
                  currentLobby.status === 'starting' ? 'bg-green-800/30 text-green-300' : 
                    'bg-blue-800/30 text-blue-300'}
              `}>
                {startCountdown ? (
                  <>
                    <Clock className="h-4 w-4 mr-1 animate-pulse" />
                    <span>Starting in {startCountdown}s</span>
                  </>
                ) : (
                  <span className="capitalize">{currentLobby.status}</span>
                )}
              </div>
            </div>

            <div className="bg-background-dark p-4 rounded-lg mb-4">
              <h3 className="text-sm uppercase tracking-wider text-gray-400 mb-3">Players {currentLobby.players.length}/{currentLobby.maxPlayers}</h3>
              
              <div className="space-y-3">
                {currentLobby.players.map(player => (
                  <div key={player.id} className="flex items-center justify-between bg-background-light p-3 rounded-lg">
                    <div className="flex items-center">
                      <div className="w-10 h-10 rounded-full bg-primary-700 flex items-center justify-center mr-3">
                        {player.username.charAt(0).toUpperCase()}
                      </div>
                      <div>
                        <div className="font-medium">{player.username}</div>
                        <div className="text-xs text-gray-400">Level {player.level}</div>
                      </div>
                    </div>
                    <div className={`
                      px-2 py-1 rounded-full text-xs
                      ${player.isReady ? 'bg-green-900/30 text-green-400' : 'bg-yellow-900/30 text-yellow-400'}
                    `}>
                      {player.isReady ? 'Ready' : 'Not Ready'}
                    </div>
                  </div>
                ))}
                
                {/* Empty slots */}
                {Array.from({ length: currentLobby.maxPlayers - currentLobby.players.length }).map((_, index) => (
                  <div key={`empty-${index}`} className="flex items-center justify-between bg-background-light/50 p-3 rounded-lg border border-dashed border-gray-700">
                    <div className="flex items-center">
                      <div className="w-10 h-10 rounded-full bg-gray-800 flex items-center justify-center mr-3">
                        ?
                      </div>
                      <div className="text-gray-500">Waiting for player...</div>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div className="flex flex-col sm:flex-row gap-3">
              <Button
                variant={isPlayerReady ? "outline" : "primary"}
                onClick={() => setPlayerReady(!isPlayerReady)}
                fullWidth
              >
                {isPlayerReady ? 'Cancel Ready' : 'Ready Up'}
              </Button>
              <Button
                variant="accent"
                onClick={leaveLobby}
                fullWidth
              >
                Leave Lobby
              </Button>
            </div>
          </div>
        </Card>

        <div className="bg-gradient-to-r from-primary-900/20 to-secondary-900/20 rounded-lg p-4">
          <h3 className="font-medium mb-2">How to start a game:</h3>
          <ul className="list-disc list-inside text-sm text-gray-300 space-y-1">
            <li>All players must click the "Ready Up" button</li>
            <li>Game will start automatically when all players are ready</li>
            <li>For 2+ player games, the minimum required players must be present</li>
            <li>You can cancel your ready status by clicking the button again</li>
          </ul>
        </div>
      </div>
    </Layout>
  );
};

export default LobbyPage;