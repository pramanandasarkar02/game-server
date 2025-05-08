import React from 'react';
import { Trophy } from 'lucide-react';
import { Game } from '../../types';
import Card from '../ui/Card';

interface GameCardProps {
  game: Game;
  onClick: () => void;
}

const GameCard: React.FC<GameCardProps> = ({ game, onClick }) => {
  const { name, description, image, difficulty } = game;

  const difficultyColors = {
    easy: 'bg-green-600',
    medium: 'bg-yellow-600',
    hard: 'bg-red-600'
  };

  return (
    <Card 
      animated
      onClick={onClick}
      className="h-full"
    >
      <div className="relative h-48 overflow-hidden">
        <img
          src={image}
          alt={name}
          className="w-full h-full object-cover transition-transform duration-500 hover:scale-110"
        />
        <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 to-transparent p-3">
          <div className="flex items-center justify-between">
            <h3 className="text-xl font-bold text-white">{name}</h3>
            <span className={`text-xs px-2 py-1 rounded-full text-white ${difficultyColors[difficulty]}`}>
              {difficulty}
            </span>
          </div>
        </div>
      </div>
      <div className="p-4">
        <p className="text-gray-300 mb-4">{description}</p>
        <div className="flex items-center text-primary-400">
          <Trophy className="h-4 w-4 mr-1" />
          <span className="text-sm">Play to earn points</span>
        </div>
      </div>
    </Card>
  );
};

export default GameCard;