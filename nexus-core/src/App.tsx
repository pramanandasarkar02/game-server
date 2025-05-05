import { useState } from 'react';

function App() {
  const [board, setBoard] = useState(Array(9).fill(''));
  const [isXNext, setIsXNext] = useState(true);
  const [winner, setWinner] = useState<string | null>(null);

  const calculateWinner = (squares: string[]) => {
    const lines = [
      [0, 1, 2], [3, 4, 5], [6, 7, 8], 
      [0, 3, 6], [1, 4, 7], [2, 5, 8], 
      [0, 4, 8], [2, 4, 6]             
    ];

    for (let line of lines) {
      const [a, b, c] = line;
      if (squares[a] && squares[a] === squares[b] && squares[a] === squares[c]) {
        return squares[a];
      }
    }
    return null;
  };

  const handleClick = (index: number) => {
    if (board[index] || winner) return;

    const newBoard = [...board];
    newBoard[index] = isXNext ? 'X' : 'O';
    setBoard(newBoard);
    
    const gameWinner = calculateWinner(newBoard);
    if (gameWinner) {
      setWinner(gameWinner);
    } else {
      setIsXNext(!isXNext);
    }
  };

  const resetGame = () => {
    setBoard(Array(9).fill(''));
    setIsXNext(true);
    setWinner(null);
  };

  const renderSquare = (index: number) => {
    return (
      <button
        className="w-16 h-16 border border-gray-400 text-2xl font-bold"
        onClick={() => handleClick(index)}
        disabled={!!board[index] || !!winner}
      >
        {board[index]}
      </button>
    );
  };

  const status = winner 
    ? `Winner: ${winner}`
    : board.every(square => square) 
    ? 'Game Draw!'
    : `Next player: ${isXNext ? 'X' : 'O'}`;

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-100">
      <h1 className="text-3xl font-bold mb-4">Tic Tac Toe</h1>
      
      <div className="mb-4 text-xl">{status}</div>
      
      <div className="grid grid-cols-3 gap-1 bg-gray-100 p-1">
        {[0, 1, 2, 3, 4, 5, 6, 7, 8].map((index) => (
          <div key={index}>{renderSquare(index)}</div>
        ))}
      </div>
      
      <button
        className="mt-6 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        onClick={resetGame}
      >
        Reset Game
      </button>
    </div>
  );
}

export default App;