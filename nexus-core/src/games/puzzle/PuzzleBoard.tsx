import React, { useState, useEffect } from 'react';
import Button from '../../components/ui/Button';

interface PuzzleBoardProps {
  size?: number;
  onGameEnd?: (moves: number) => void;
}

const PuzzleBoard: React.FC<PuzzleBoardProps> = ({ 
  size = 3,
  onGameEnd 
}) => {
  const [tiles, setTiles] = useState<number[]>([]);
  const [emptyIndex, setEmptyIndex] = useState<number>(size * size - 1);
  const [moves, setMoves] = useState<number>(0);
  const [isComplete, setIsComplete] = useState<boolean>(false);
  const [timerSeconds, setTimerSeconds] = useState<number>(0);
  const [isTimerRunning, setIsTimerRunning] = useState<boolean>(false);

  // Initialize board
  const initializeBoard = () => {
    // Create solved board
    const initialTiles = Array.from({ length: size * size }, (_, i) => i);
    // Shuffle
    const shuffledTiles = [...initialTiles];
    
    // Shuffle until we get a solvable puzzle
    let isSolvable = false;
    while (!isSolvable) {
      for (let i = shuffledTiles.length - 2; i > 0; i--) {
        const j = Math.floor(Math.random() * i);
        [shuffledTiles[i], shuffledTiles[j]] = [shuffledTiles[j], shuffledTiles[i]];
      }
      
      // Check if puzzle is solvable
      // A puzzle is solvable if the number of inversions is even
      let inversions = 0;
      for (let i = 0; i < shuffledTiles.length - 1; i++) {
        if (shuffledTiles[i] === 0) continue;
        for (let j = i + 1; j < shuffledTiles.length; j++) {
          if (shuffledTiles[j] === 0) continue;
          if (shuffledTiles[i] > shuffledTiles[j]) {
            inversions++;
          }
        }
      }
      
      // For odd-sized boards, the puzzle is solvable if inversions is even
      isSolvable = inversions % 2 === 0;
    }
    
    // Find empty tile position
    const emptyPos = shuffledTiles.indexOf(0);
    
    setTiles(shuffledTiles);
    setEmptyIndex(emptyPos);
    setMoves(0);
    setIsComplete(false);
    setTimerSeconds(0);
    setIsTimerRunning(true);
  };

  // Check if puzzle is complete
  const checkCompletion = (currentTiles: number[]) => {
    for (let i = 0; i < currentTiles.length - 1; i++) {
      if (currentTiles[i] !== i + 1) {
        return false;
      }
    }
    return currentTiles[currentTiles.length - 1] === 0;
  };

  // Handle tile click
  const handleTileClick = (index: number) => {
    if (isComplete) return;
    
    // Check if tile is adjacent to empty space
    const row = Math.floor(index / size);
    const col = index % size;
    const emptyRow = Math.floor(emptyIndex / size);
    const emptyCol = emptyIndex % size;
    
    const isAdjacent = 
      (row === emptyRow && Math.abs(col - emptyCol) === 1) || 
      (col === emptyCol && Math.abs(row - emptyRow) === 1);
    
    if (isAdjacent) {
      // Swap tiles
      const newTiles = [...tiles];
      [newTiles[index], newTiles[emptyIndex]] = [newTiles[emptyIndex], newTiles[index]];
      
      setTiles(newTiles);
      setEmptyIndex(index);
      setMoves(moves + 1);
      
      // Check if puzzle is complete
      if (checkCompletion(newTiles)) {
        setIsComplete(true);
        setIsTimerRunning(false);
        if (onGameEnd) onGameEnd(moves + 1);
      }
    }
  };

  // Timer effect
  useEffect(() => {
    let timer: NodeJS.Timeout | null = null;
    
    if (isTimerRunning) {
      timer = setInterval(() => {
        setTimerSeconds(prev => prev + 1);
      }, 1000);
    }
    
    return () => {
      if (timer) clearInterval(timer);
    };
  }, [isTimerRunning]);

  // Initialize on component mount
  useEffect(() => {
    initializeBoard();
  }, []);

  // Format time (mm:ss)
  const formatTime = (totalSeconds: number) => {
    const minutes = Math.floor(totalSeconds / 60);
    const seconds = totalSeconds % 60;
    return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
  };

  return (
    <div className="w-full max-w-md mx-auto">
      <div className="flex justify-between items-center mb-4">
        <div className="bg-primary-900/30 px-4 py-2 rounded-lg">
          <span className="text-primary-300">Moves: {moves}</span>
        </div>
        <div className="bg-accent-900/30 px-4 py-2 rounded-lg">
          <span className="text-accent-300">Time: {formatTime(timerSeconds)}</span>
        </div>
      </div>

      <div className="mb-4 text-center py-2 bg-background-light rounded-lg">
        <div className="text-lg font-semibold">
          {isComplete ? 'Puzzle Completed!' : 'Arrange in order'}
        </div>
      </div>

      <div 
        className="grid gap-2 mb-4 aspect-square" 
        style={{ gridTemplateColumns: `repeat(${size}, 1fr)` }}
      >
        {tiles.map((tile, index) => (
          <button
            key={index}
            className={`
              aspect-square flex items-center justify-center text-2xl font-bold
              ${tile === 0 ? 'bg-transparent cursor-default' : 'bg-background-light hover:bg-primary-800/20 cursor-pointer'}
              ${tile === index + 1 && tile !== 0 ? 'border-2 border-green-500' : 'border border-gray-700'}
              rounded-md transition-colors duration-200
            `}
            onClick={() => tile !== 0 && handleTileClick(index)}
            disabled={isComplete || tile === 0}
          >
            {tile !== 0 && tile}
          </button>
        ))}
      </div>

      <Button onClick={initializeBoard} fullWidth>
        New Game
      </Button>
    </div>
  );
};

export default PuzzleBoard;