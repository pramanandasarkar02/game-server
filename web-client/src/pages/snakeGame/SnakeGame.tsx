import React, { useEffect, useRef, useState } from "react";
import { useParams } from "react-router-dom";

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
  const [chatMessage, setChatMessage] = useState("");
  const [chatLog, setChatLog] = useState<string[]>([]);
  const wsRef = useRef<WebSocket | null>(null);
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const { gameId, userId } = useParams();

  // connect WebSocket
  useEffect(() => {
    const ws = new WebSocket(
      `ws://localhost:8080/ws?matchId=${gameId}&playerId=${userId}`
    );
    wsRef.current = ws;

    ws.onopen = () => console.log("Connected to WebSocket");
    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);

        if (data.type === "chat") {
          setChatLog((prev) => [...prev, `${data.from}: ${data.message}`]);
        } else if (data.type === "update") {
          setGameState(data.state ?? data); // handle if server wraps inside {type, state}
        }
      } catch (err) {
        console.error("Invalid JSON:", event.data);
      }
    };
    ws.onerror = (err) => console.error("WebSocket error:", err);
    ws.onclose = () => console.log("Disconnected");

    return () => ws.close();
  }, [gameId, userId]);

  // draw game
  useEffect(() => {
    if (!gameState || !canvasRef.current) return;
    const ctx = canvasRef.current.getContext("2d");
    if (!ctx) return;
    ctx.clearRect(0, 0, BOARD_WIDTH * CELL_SIZE, BOARD_HEIGHT * CELL_SIZE);

    // Player snake
    const snake = gameState.playerSnake;
    if (snake) {
      ctx.fillStyle = "green";
      snake.snakeBody.forEach((part) =>
        ctx.fillRect(
          part.x * CELL_SIZE,
          part.y * CELL_SIZE,
          CELL_SIZE,
          CELL_SIZE
        )
      );
      ctx.fillStyle = "darkgreen";
      ctx.fillRect(
        snake.snakeHead.x * CELL_SIZE,
        snake.snakeHead.y * CELL_SIZE,
        CELL_SIZE,
        CELL_SIZE
      );
    }

    // Other snakes
    ctx.fillStyle = "yellow";
    gameState.otherSnakes?.forEach((other) => {
      other.snake.snakeBody.forEach((part) =>
        ctx.fillRect(
          part.x * CELL_SIZE,
          part.y * CELL_SIZE,
          CELL_SIZE,
          CELL_SIZE
        )
      );
      ctx.fillRect(
        other.snake.snakeHead.x * CELL_SIZE,
        other.snake.snakeHead.y * CELL_SIZE,
        CELL_SIZE,
        CELL_SIZE
      );
    });

    // Foods
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

    // Obstacles
    ctx.fillStyle = "gray";
    gameState.obstacles.forEach((obs) =>
      obs.object?.forEach((p) =>
        ctx.fillRect(p.x * CELL_SIZE, p.y * CELL_SIZE, CELL_SIZE, CELL_SIZE)
      )
    );
  }, [gameState]);

  // send direction change
  useEffect(() => {
    const handleKey = (e: KeyboardEvent) => {
      const keyMap: Record<string, string> = {
        ArrowUp: "UP",
        ArrowDown: "DOWN",
        ArrowLeft: "LEFT",
        ArrowRight: "RIGHT",
        w: "UP",
        s: "DOWN",
        a: "LEFT",
        d: "RIGHT",
      };
      const dir = keyMap[e.key];
      if (dir && wsRef.current?.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify({ type: "move", direction: dir }));
      }
    };
    window.addEventListener("keydown", handleKey);
    return () => window.removeEventListener("keydown", handleKey);
  }, []);

  // send chat
  const sendChat = () => {
    if (!chatMessage.trim() || wsRef.current?.readyState !== WebSocket.OPEN)
      return;
    wsRef.current.send(
      JSON.stringify({
        type: "chat",
        from: userId,
        message: chatMessage,
      })
    );
    setChatMessage("");
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-black text-white">
      <h1 className="text-2xl font-bold mb-4">Snake Game</h1>
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

      {/* Chat UI */}
      <div className="mt-6 w-1/3">
        <div className="h-40 overflow-y-auto bg-gray-900 p-2 rounded text-sm border border-gray-700 mb-2">
          {chatLog.map((msg, idx) => (
            <div key={idx}>{msg}</div>
          ))}
        </div>
        <div className="flex">
          <input
            className="flex-1 p-2 bg-gray-800 border border-gray-600 rounded-l text-white"
            value={chatMessage}
            onChange={(e) => setChatMessage(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && sendChat()}
            placeholder="Type a message..."
          />
          <button className="px-4 bg-blue-600 rounded-r" onClick={sendChat}>
            Send
          </button>
        </div>
      </div>
    </div>
  );
};

export default SnakeGame;
