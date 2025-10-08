import { useCallback, useEffect, useState } from "react";
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

const GRID_SIZE = 20;
const CELL_SIZE = 20;
const INITIAL_SNAKE: Position[] = [
  { x: 10, y: 10 },
  { x: 9, y: 10 },
  { x: 8, y: 10 },
];

const INITIAL_DIRECTION: Direction = "RIGHT";
const GAME_SPEED = 150; // smoother speed

const SnakeGame = () => {
  const { gameId, playerId } = useParams<SnakeGameParams>();

  const [snake, setSnake] = useState<Position[]>(INITIAL_SNAKE);
  const [direction, setDirection] = useState<Direction>(INITIAL_DIRECTION);
  const [nextDirection, setNextDirection] = useState<Direction>(INITIAL_DIRECTION);
  const [gameOver, setGameOver] = useState<boolean>(false);
  const [score, setScore] = useState<number>(0);
  const [food, setFood] = useState<Position>(() => ({
    x: Math.floor(Math.random() * GRID_SIZE),
    y: Math.floor(Math.random() * GRID_SIZE),
  }));

  // Move snake
  const moveSnake = useCallback(() => {
    if (gameOver) return;

    setDirection(nextDirection);

    setSnake((prevSnake) => {
      const head = prevSnake[0];
      let newHead: Position = { ...head };

      switch (nextDirection) {
        case "UP":
          newHead = { x: head.x, y: head.y - 1 };
          break;
        case "DOWN":
          newHead = { x: head.x, y: head.y + 1 };
          break;
        case "LEFT":
          newHead = { x: head.x - 1, y: head.y };
          break;
        case "RIGHT":
          newHead = { x: head.x + 1, y: head.y };
          break;
      }

      // border collision
      if (
        newHead.x < 0 ||
        newHead.x >= GRID_SIZE ||
        newHead.y < 0 ||
        newHead.y >= GRID_SIZE
      ) {
        setGameOver(true);
        return prevSnake;
      }

      // self collision
      if (prevSnake.some((seg) => seg.x === newHead.x && seg.y === newHead.y)) {
        setGameOver(true);
        return prevSnake;
      }

      const newSnake = [newHead, ...prevSnake];

      // food collision
      if (newHead.x === food.x && newHead.y === food.y) {
        setScore((s) => s + 1);
        setFood({
          x: Math.floor(Math.random() * GRID_SIZE),
          y: Math.floor(Math.random() * GRID_SIZE),
        });
      } else {
        newSnake.pop(); // move forward (no grow)
      }

      return newSnake;
    });
  }, [nextDirection, gameOver, food]);

  // handle key press
  useEffect(() => {
    const handleKeyPress = (e: KeyboardEvent) => {
      switch (e.key) {
        case "ArrowUp":
        case "w":
          setNextDirection((prev) => (prev !== "DOWN" ? "UP" : prev));
          break;
        case "ArrowDown":
        case "s":
          setNextDirection((prev) => (prev !== "UP" ? "DOWN" : prev));
          break;
        case "ArrowLeft":
        case "a":
          setNextDirection((prev) => (prev !== "RIGHT" ? "LEFT" : prev));
          break;
        case "ArrowRight":
        case "d":
          setNextDirection((prev) => (prev !== "LEFT" ? "RIGHT" : prev));
          break;
      }
    };
    window.addEventListener("keydown", handleKeyPress);
    return () => window.removeEventListener("keydown", handleKeyPress);
  }, []);

  // game loop
  useEffect(() => {
    const interval = setInterval(moveSnake, GAME_SPEED);
    return () => clearInterval(interval);
  }, [moveSnake]);

  const restartGame = () => {
    setSnake(INITIAL_SNAKE);
    setDirection(INITIAL_DIRECTION);
    setNextDirection(INITIAL_DIRECTION);
    setGameOver(false);
    setScore(0);
    setFood({
      x: Math.floor(Math.random() * GRID_SIZE),
      y: Math.floor(Math.random() * GRID_SIZE),
    });
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-4 bg-gray-900">
      <div className="bg-gray-800 rounded-lg px-4 py-3 shadow-lg">
        <div className="text-center mb-4">
          <h1 className="text-2xl font-bold text-gray-100">Snake Game</h1>
          <div className="flex justify-between items-center text-sm text-gray-400 mb-1">
            <span>Game: {gameId}</span>
            <span>Player: {playerId}</span>
          </div>
          <div className="text-3xl font-bold text-green-400">Score: {score}</div>
          {gameOver && (
            <div className="text-red-400 font-semibold mt-2">
              Game Over!
              <button
                onClick={restartGame}
                className="ml-3 bg-blue-500 px-3 py-1 rounded hover:bg-blue-600 text-white"
              >
                Restart
              </button>
            </div>
          )}
        </div>

        <div
          className="relative bg-gray-700 rounded overflow-hidden shadow-lg"
          style={{
            width: `${GRID_SIZE * CELL_SIZE}px`,
            height: `${GRID_SIZE * CELL_SIZE}px`,
            display: "grid",
            gridTemplateColumns: `repeat(${GRID_SIZE}, ${CELL_SIZE}px)`,
          }}
        >
          {Array.from({ length: GRID_SIZE * GRID_SIZE }).map((_, idx) => {
            const x = idx % GRID_SIZE;
            const y = Math.floor(idx / GRID_SIZE);

            const isSnakeHead = snake[0].x === x && snake[0].y === y;
            const isSnakeBody = snake.some(
              (seg, i) => i !== 0 && seg.x === x && seg.y === y
            );
            const isFood = food.x === x && food.y === y;

            return (
              <div
                key={idx}
                style={{
                  width: `${CELL_SIZE}px`,
                  height: `${CELL_SIZE}px`,
                }}
                className={`${
                  isSnakeHead
                    ? "bg-yellow-400"
                    : isSnakeBody
                    ? "bg-green-400"
                    : isFood
                    ? "bg-red-500"
                    : "bg-gray-600"
                }`}
              />
            );
          })}
        </div>
      </div>
    </div>
  );
};

export default SnakeGame;
