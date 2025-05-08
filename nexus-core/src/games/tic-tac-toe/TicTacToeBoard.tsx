import React, { useState, useEffect } from 'react';
import Button from '../../components/ui/Button';

interface TicTacToeBoardProps {
  opponent?: string;
  onGameEnd?: (result: 'win' | 'lose' | 'draw') => void;
}

const TicTacToeBoard: React.FC<TicTacToeBoardProps> = ({ 
  opponent = 'Computer', 
  onGameEnd 
}) => {
  const [board, setBoard] = useState<Array<string | null>>(Array(9).fill(null));
  const [isXNext, setIsXNext] = useState(true);
  const [winner, setWinner] = useState<string | null>(null);
  const [isDraw, setIsDraw] = useState(false);
  const [score, setScore] = useState({ player: 0, opponent: 0 });

  const calculateWinner = (squares: Array<string | null>): string | null => {
    const lines = [
      [0, 1, 2],
      [3, 4, 5],
      [6, 7, 8],
      [0, 3, 6],
      [1, 4, 7],
      [2, 5, 8],
      [0, 4, 8],
      [2, 4, 6],
    ];

    for (let i = 0; i < lines.length; i++) {
      const [a, b, c] = lines[i];
      if (squares[a] && squares[a] === squares[b] && squares[a] === squares[c]) {
        return squares[a]!;
      }
    }
    return null;
  };

  const makeComputerMove = () => {
    // Simple AI - finds first empty cell
    const newBoard = [...board];
    const emptySquares = newBoard
      .map((square, index) => (square === null ? index : null))
      .filter((index) => index !== null) as number[];

    if (emptySquares.length > 0) {
      // Choose a random empty square
      const randomIndex = Math.floor(Math.random() * emptySquares.length);
      const computerMove = emptySquares[randomIndex];
      newBoard[computerMove] = 'O';
      setBoard(newBoard);
      setIsXNext(true);
      
      const computerWinner = calculateWinner(newBoard);
      if (computerWinner) {
        setWinner(computerWinner);
        setScore(prev => ({ ...prev, opponent: prev.opponent + 1 }));
        if (onGameEnd) onGameEnd('lose');
      } else if (!newBoard.includes(null)) {
        setIsDraw(true);
        if (onGameEnd) onGameEnd('draw');
      }
    }
  };

  useEffect(() => {
    // If it's computer's turn, make a move
    if (!isXNext && !winner && !isDraw && opponent === 'Computer') {
      const timer = setTimeout(() => {
        makeComputerMove();
      }, 500);
      return () => clearTimeout(timer);
    }
  }, [isXNext, winner, isDraw]);

  const handleClick = (index: number) => {
    if (board[index] || winner || !isXNext) return;

    const newBoard = [...board];
    newBoard[index] = 'X';
    setBoard(newBoard);
    setIsXNext(false);

    const newWinner = calculateWinner(newBoard);
    if (newWinner) {
      setWinner(newWinner);
      setScore(prev => ({ ...prev, player: prev.player + 1 }));
      if (onGameEnd) onGameEnd('win');
    } else if (!newBoard.includes(null)) {
      setIsDraw(true);
      if (onGameEnd) onGameEnd('draw');
    }
  };

  const resetGame = () => {
    setBoard(Array(9).fill(null));
    setIsXNext(true);
    setWinner(null);
    setIsDraw(false);
  };

  const renderSquare = (index: number) => {
    return (
      <button
        className={`
          w-full h-20 bg-background-light border border-gray-700 text-3xl font-bold
          hover:bg-primary-800/20 focus:outline-none transition-colors duration-200
          ${board[index] === 'X' ? 'text-primary-500' : 'text-accent-500'}
        `}
        onClick={() => handleClick(index)}
        disabled={!!board[index] || !!winner || isDraw || !isXNext}
      >
        {board[index]}
      </button>
    );
  };

  const getStatus = () => {
    if (winner) {
      return `Winner: ${winner === 'X' ? 'You' : opponent}`;
    } else if (isDraw) {
      return 'Game Ended in Draw';
    } else {
      return `Next player: ${isXNext ? 'You (X)' : `${opponent} (O)`}`;
    }
  };

  return (
    <div className="w-full max-w-md mx-auto">
      <div className="flex justify-between items-center mb-4">
        <div className="bg-primary-900/30 px-4 py-2 rounded-lg">
          <span className="text-primary-300">You: {score.player}</span>
        </div>
        <div className="bg-accent-900/30 px-4 py-2 rounded-lg">
          <span className="text-accent-300">{opponent}: {score.opponent}</span>
        </div>
      </div>

      <div className="mb-4 text-center py-2 bg-background-light rounded-lg">
        <div className="text-lg font-semibold">{getStatus()}</div>
      </div>

      <div className="grid grid-cols-3 gap-2 mb-4">
        {[0, 1, 2, 3, 4, 5, 6, 7, 8].map((index) => (
          <div key={index}>{renderSquare(index)}</div>
        ))}
      </div>

      <Button onClick={resetGame} fullWidth>
        New Game
      </Button>
    </div>
  );
};

export default TicTacToeBoard;