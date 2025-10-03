import { useCallback, useState } from "react";
import { useParams } from "react-router-dom";

type SnakeGameParams = {
  gameId: string;
  playerId: string;
};

type Position = {
  x: number;
  y: number;
};

type Direction = "UP" | "DOWN" | "LEFT" | "RIGHT";

const GRID_SIZE = 40;
const CELL_SIZE = 20;
const INITIAL_SNAKE: Position[] = [
  { x: GRID_SIZE / 2, y: CELL_SIZE / 2 },
  { x: GRID_SIZE / 2 - 1, y: CELL_SIZE / 2 },
  { x: GRID_SIZE / 2 - 2, y: CELL_SIZE / 2 },
];

const INITIAL_DIRECTION: Direction = "RIGHT";
const GAME_SPEED = 150;

const SnakeGame = () => {
  const { gameId, playerId } = useParams<SnakeGameParams>();

  const [snake, setSnake] = useState<Position[]>(INITIAL_SNAKE);
  const [direction, setDirection] = useState<Direction>(INITIAL_DIRECTION);
  const [gameOver, setGameOver] = useState<boolean>(false);
  const [nextDirection, setNextDirection] = useState<Direction>(INITIAL_DIRECTION)

  const moveSnake = useCallback(()=> {
    if (gameOver ){
        return 
    }
    setDirection(nextDirection);

    // setSnake(()=> {
    //     const head = pr
    // })

  }, [])

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <div className="bg-gray-800 rounded-lg max-w-2xl px-4 py-2">
        <div className="text-center mb-4">
          <h1 className="text-2xl font-bold text-gray-200">SnakeGame</h1>
          <div className="flex justify-between items-center text-sm text-gray-400 mb-2">
            <span>Game: {gameId}</span>
            <span>Player: {playerId}</span>
          </div>
          <div className="text-3xl font-bold text-green-500">Score: 0</div>
        </div>
        <div className="relative bg-gray-300 rounded-lg overflow-hidden border-4 border-gray-400 shadow-lg">
          <div className="grid gap-0"
          style={{
            gridTemplateColumns: `repeat(${GRID_SIZE}, ${CELL_SIZE}px)`,
            width: `${GRID_SIZE * CELL_SIZE}px`,
            height: `${GRID_SIZE / 2 * CELL_SIZE}px`
          }}>
            {Array.from({ length: GRID_SIZE * GRID_SIZE}).map((_, idx) => {
                const x = idx % GRID_SIZE;
                const y = Math.floor(idx / GRID_SIZE);
                const isSnakeHead = snake[0].x === x && snake[0].y === y;
                const isSnakeBody = snake.some((segment) => segment.x === x && segment.y === y);


                return (
                    <div
                    key={idx}
                    className={`${
                        isSnakeHead?"bg-amber-500":
                        isSnakeBody?"bg-green-400":"bg-gray-500"
                    }`}
                    style={{
                        width: `${CELL_SIZE}px`,
                        height: `${CELL_SIZE}px`
                    }}
                    />
                

                    
                )
            })}

          </div>
        </div>
      </div>
    </div>
  );
};

export default SnakeGame;
