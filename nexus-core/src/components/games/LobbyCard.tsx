import React from 'react';
import { Users } from 'lucide-react';
import { Lobby } from '../../types';
import Card from '../ui/Card';
import Button from '../ui/Button';

interface LobbyCardProps {
  lobby: Lobby;
  onJoin: () => void;
}

const LobbyCard: React.FC<LobbyCardProps> = ({ lobby, onJoin }) => {
  const { gameName, players, maxPlayers, status } = lobby;
  const isFull = players.length >= maxPlayers;

  return (
    <Card className="mb-4">
      <div className="p-4">
        <div className="flex justify-between items-center mb-2">
          <h3 className="text-lg font-semibold">{gameName}</h3>
          <span className={`
            px-2 py-1 text-xs rounded-full 
            ${status === 'waiting' ? 'bg-yellow-700 text-yellow-100' : 
              status === 'starting' ? 'bg-green-700 text-green-100' : 
                'bg-blue-700 text-blue-100'}
          `}>
            {status === 'waiting' ? 'Waiting' : 
             status === 'starting' ? 'Starting' : 'Playing'}
          </span>
        </div>
        
        <div className="flex items-center text-gray-400 mb-3">
          <Users className="h-4 w-4 mr-1" />
          <span className="text-sm">
            {players.length}/{maxPlayers} players
          </span>
        </div>

        <div className="mb-4">
          <h4 className="text-sm font-medium text-gray-300 mb-2">Players:</h4>
          <div className="space-y-2">
            {players.map(player => (
              <div 
                key={player.id}
                className="flex items-center justify-between bg-background-dark rounded p-2"
              >
                <div className="flex items-center">
                  <div className="w-8 h-8 rounded-full bg-primary-700 flex items-center justify-center mr-2">
                    {player.username.charAt(0).toUpperCase()}
                  </div>
                  <span>{player.username}</span>
                </div>
                <div className="flex items-center">
                  <span className="text-xs mr-2">Lv.{player.level}</span>
                  <span className={`w-2 h-2 rounded-full ${player.isReady ? 'bg-green-500' : 'bg-yellow-500'}`}></span>
                </div>
              </div>
            ))}
          </div>
        </div>

        <Button 
          onClick={onJoin}
          variant={isFull ? 'outline' : 'primary'}
          disabled={isFull}
          fullWidth
        >
          {isFull ? 'Lobby Full' : 'Join Lobby'}
        </Button>
      </div>
    </Card>
  );
};

export default LobbyCard;