import React, { useEffect, useRef, useState } from "react";

const CELL_SIZE = 10;
const BOARD_WIDTH = 60;
const BOARD_HEIGHT = 40;

interface Point {
  x: number;
  y: number;
}

interface Score {
  value: number;
}

interface Snake {
  snakeHead: Point;
  snakeBody: Point[];
  direction: string;
  score: Score;
  time: string;
}

interface Food {
  position: Point;
  value: number;
}

interface Obstacle {
  object: Point[];
}

interface GameState {
  playerId: string;
  playerSnake: Snake;
  foods: Food[];
  otherSnakes: { playerId: string; snake: Snake }[];
  obstacles: Obstacle[];
}

const SnakeGame: React.FC = () => {
  const [gameState, setGameState] = useState<GameState | null>(null);
  const canvasRef = useRef<HTMLCanvasElement | null>(null);

  // WebSocket connection
  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/ws?matchId=abc123&playerId=p1");

    ws.onopen = () => console.log("✅ Connected to WebSocket");
    ws.onmessage = (event) => {
      try {
        const data: GameState = JSON.parse(event.data);
        setGameState(data);
      } catch (err) {
        console.error("Invalid JSON:", err);
      }
    };
    ws.onerror = (err) => console.error("WebSocket error:", err);
    ws.onclose = () => console.log("❌ Disconnected from WebSocket");

    return () => ws.close();
  }, []);

  // Draw game board
  useEffect(() => {
    if (!gameState || !canvasRef.current) return;
    const ctx = canvasRef.current.getContext("2d");
    if (!ctx) return;

    ctx.clearRect(0, 0, BOARD_WIDTH * CELL_SIZE, BOARD_HEIGHT * CELL_SIZE);

    // Draw player's snake
    const snake = gameState.playerSnake;
    if (snake) {
      ctx.fillStyle = "green";
      snake.snakeBody.forEach((part) => {
        ctx.fillRect(part.x * CELL_SIZE, part.y * CELL_SIZE, CELL_SIZE, CELL_SIZE);
      });

      // Draw snake head
      ctx.fillStyle = "darkgreen";
      ctx.fillRect(
        snake.snakeHead.x * CELL_SIZE,
        snake.snakeHead.y * CELL_SIZE,
        CELL_SIZE,
        CELL_SIZE
      );
    }

    // Draw other snakes
    ctx.fillStyle = "yellow";
    gameState.otherSnakes?.forEach((other) => {
      other.snake.snakeBody.forEach((part) => {
        ctx.fillRect(part.x * CELL_SIZE, part.y * CELL_SIZE, CELL_SIZE, CELL_SIZE);
      });
      ctx.fillRect(
        other.snake.snakeHead.x * CELL_SIZE,
        other.snake.snakeHead.y * CELL_SIZE,
        CELL_SIZE,
        CELL_SIZE
      );
    });

    // Draw foods
    ctx.fillStyle = "red";
    gameState.foods.forEach((food) => {
      ctx.beginPath();
      ctx.arc(
        food.position.x * CELL_SIZE + CELL_SIZE / 2,
        food.position.y * CELL_SIZE + CELL_SIZE / 2,
        CELL_SIZE / 2,
        0,
        Math.PI * 2
      );
      ctx.fill();
    });

    // Draw obstacles
    ctx.fillStyle = "gray";
    gameState.obstacles.forEach((obs) => {
      obs.object?.forEach((p) => {
        ctx.fillRect(p.x * CELL_SIZE, p.y * CELL_SIZE, CELL_SIZE, CELL_SIZE);
      });
    });
  }, [gameState]);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-black text-white">
      <h1 className="text-xl font-bold mb-4">🐍 Snake Game</h1>
      <canvas
        ref={canvasRef}
        width={BOARD_WIDTH * CELL_SIZE}
        height={BOARD_HEIGHT * CELL_SIZE}
        style={{ border: "2px solid white", background: "#111" }}
      />
      {gameState && (
        <div className="mt-4 text-center">
          <p>Player: {gameState.playerId}</p>
          <p>Score: {gameState.playerSnake?.score?.value ?? 0}</p>
        </div>
      )}
    </div>
  );
};

export default SnakeGame;
