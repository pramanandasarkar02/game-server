import React, { useEffect, useRef, useState, useCallback } from "react";
import { useParams } from "react-router-dom";
// import PlayerContext from "../../context/PlayerContext";

const CELL_SIZE = 16;
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
  otherSnakes: Snake[];
  obstacles: Obstacle[];
}

interface ChatMessage {
  type: string;
  from: string;
  message: string;
}

const SnakeGame: React.FC = () => {
  const [gameState, setGameState] = useState<GameState | null>(null);
  const [chatMessage, setChatMessage] = useState("");
  const [chatLog, setChatLog] = useState<ChatMessage[]>([]);
  const [connectionStatus, setConnectionStatus] = useState<"connecting" | "connected" | "disconnected">("connecting");
  const [error, setError] = useState<string | null>(null);
  
  const wsRef = useRef<WebSocket | null>(null);
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const lastDirectionRef = useRef<string>("");
  
  const { gameId, userId } = useParams<{ gameId: string; userId: string }>();
  // const {player} = useContext(PlayerContext);
  // const username = player?.username;
  // Draw game state on canvas
  const drawGame = useCallback(() => {
    if (!gameState || !canvasRef.current) return;
    
    const ctx = canvasRef.current.getContext("2d");
    if (!ctx) return;

    // Clear canvas
    ctx.fillStyle = "#111";
    ctx.fillRect(0, 0, BOARD_WIDTH * CELL_SIZE, BOARD_HEIGHT * CELL_SIZE);

    // Draw grid (optional)
    ctx.strokeStyle = "#222";
    ctx.lineWidth = 0.5;
    for (let x = 0; x <= BOARD_WIDTH; x++) {
      ctx.beginPath();
      ctx.moveTo(x * CELL_SIZE, 0);
      ctx.lineTo(x * CELL_SIZE, BOARD_HEIGHT * CELL_SIZE);
      ctx.stroke();
    }
    for (let y = 0; y <= BOARD_HEIGHT; y++) {
      ctx.beginPath();
      ctx.moveTo(0, y * CELL_SIZE);
      ctx.lineTo(BOARD_WIDTH * CELL_SIZE, y * CELL_SIZE);
      ctx.stroke();
    }

    // Draw obstacles
    if (gameState.obstacles && gameState.obstacles.length > 0) {
      ctx.fillStyle = "#555";
      gameState.obstacles.forEach((obs) => {
        if (obs && obs.object && obs.object.length > 0) {
          obs.object.forEach((p) => {
            ctx.fillRect(
              p.x * CELL_SIZE,
              p.y * CELL_SIZE,
              CELL_SIZE,
              CELL_SIZE
            );
          });
        }
      });
    }

    // Draw foods
    if (gameState.foods && gameState.foods.length > 0) {
      gameState.foods.forEach((food) => {
        if (food && food.position) {
          // Draw food as a circle with value
          ctx.fillStyle = "#ff4444";
          ctx.beginPath();
          ctx.arc(
            food.position.x * CELL_SIZE + CELL_SIZE / 2,
            food.position.y * CELL_SIZE + CELL_SIZE / 2,
            CELL_SIZE / 2 - 1,
            0,
            Math.PI * 2
          );
          ctx.fill();
          
          // Draw food value
          ctx.fillStyle = "white";
          ctx.font = "12px Arial";
          ctx.textAlign = "center";
          ctx.textBaseline = "middle";
          ctx.fillText(
            food.value.toString(),
            food.position.x * CELL_SIZE + CELL_SIZE / 2,
            food.position.y * CELL_SIZE + CELL_SIZE / 2
          );
        }
      });
    }

    // Draw other snakes
    if (gameState.otherSnakes && gameState.otherSnakes.length > 0) {
      gameState.otherSnakes.forEach((snake, index) => {
        if (snake && snake.snakeBody && snake.snakeHead) {
          // Different colors for different snakes
          const colors = ["#ffaa00", "#00aaff", "#ff00aa", "#aaff00"];
          const color = colors[index % colors.length];
          
          // Draw body
          ctx.fillStyle = color;
          snake.snakeBody.forEach((part) => {
            ctx.fillRect(
              part.x * CELL_SIZE + 1,
              part.y * CELL_SIZE + 1,
              CELL_SIZE - 2,
              CELL_SIZE - 2
            );
          });
          
          // Draw head (slightly different)
          ctx.fillStyle = color;
          ctx.fillRect(
            snake.snakeHead.x * CELL_SIZE,
            snake.snakeHead.y * CELL_SIZE,
            CELL_SIZE,
            CELL_SIZE
          );
          
          // Draw eyes
          ctx.fillStyle = "white";
          ctx.fillRect(
            snake.snakeHead.x * CELL_SIZE + 2,
            snake.snakeHead.y * CELL_SIZE + 2,
            2,
            2
          );
          ctx.fillRect(
            snake.snakeHead.x * CELL_SIZE + 6,
            snake.snakeHead.y * CELL_SIZE + 2,
            2,
            2
          );
        }
      });
    }

    // Draw player snake (on top)
    const snake = gameState.playerSnake;
    if (snake && snake.snakeBody && snake.snakeHead) {
      // Draw body
      ctx.fillStyle = "#22cc22";
      snake.snakeBody.forEach((part) => {
        ctx.fillRect(
          part.x * CELL_SIZE + 1,
          part.y * CELL_SIZE + 1,
          CELL_SIZE - 2,
          CELL_SIZE - 2
        );
      });
      
      // Draw head
      ctx.fillStyle = "#00ff00";
      ctx.fillRect(
        snake.snakeHead.x * CELL_SIZE,
        snake.snakeHead.y * CELL_SIZE,
        CELL_SIZE,
        CELL_SIZE
      );
      
      // Draw eyes
      ctx.fillStyle = "black";
      ctx.fillRect(
        snake.snakeHead.x * CELL_SIZE + 2,
        snake.snakeHead.y * CELL_SIZE + 2,
        2,
        2
      );
      ctx.fillRect(
        snake.snakeHead.x * CELL_SIZE + 6,
        snake.snakeHead.y * CELL_SIZE + 2,
        2,
        2
      );
    }
  }, [gameState]);

  // Connect WebSocket
  const connectWebSocket = useCallback(() => {
    if (!gameId || !userId) {
      setError("Missing game ID or user ID");
      return;
    }

    try {
      const ws = new WebSocket(
        `ws://localhost:8080/ws?matchId=${gameId}&playerId=${userId}`
      );
      
      wsRef.current = ws;
      setConnectionStatus("connecting");

      ws.onopen = () => {
        console.log("Connected to WebSocket");
        setConnectionStatus("connected");
        setError(null);
      };

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          console.log("Received:", data);

          if (data.type === "chat") {
            setChatLog((prev) => [
              ...prev,
              { type: "chat", from: data.from, message: data.message },
            ]);
          } else if (data.type === "update" && data.state) {
            setGameState(data.state);
          } else if (data.playerId) {
            // Direct game state
            setGameState(data);
          }
        } catch (err) {
          console.error("Invalid JSON:", event.data, err);
        }
      };

      ws.onerror = (err) => {
        console.error("WebSocket error:", err);
        setError("WebSocket connection error");
      };

      ws.onclose = () => {
        console.log("Disconnected from WebSocket");
        setConnectionStatus("disconnected");
        
        // Attempt to reconnect after 3 seconds
        if (reconnectTimeoutRef.current) {
          clearTimeout(reconnectTimeoutRef.current);
        }
        reconnectTimeoutRef.current = setTimeout(() => {
          console.log("Attempting to reconnect...");
          connectWebSocket();
        }, 3000);
      };
    } catch (err) {
      console.error("Failed to create WebSocket:", err);
      setError("Failed to connect to game server");
      setConnectionStatus("disconnected");
    }
  }, [gameId, userId]);

  // Initialize WebSocket connection
  useEffect(() => {
    connectWebSocket();

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [connectWebSocket]);

  // Draw game on state change
  useEffect(() => {
    drawGame();
  }, [drawGame]);

  // Handle keyboard input
  useEffect(() => {
    const handleKey = (e: KeyboardEvent) => {
      // Prevent scrolling with arrow keys
      if (["ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight"].includes(e.key)) {
        e.preventDefault();
      }

      const keyMap: Record<string, string> = {
        ArrowUp: "UP",
        ArrowDown: "DOWN",
        ArrowLeft: "LEFT",
        ArrowRight: "RIGHT",
        w: "UP",
        W: "UP",
        s: "DOWN",
        S: "DOWN",
        a: "LEFT",
        A: "LEFT",
        d: "RIGHT",
        D: "RIGHT",
      };

      const dir = keyMap[e.key];
      
      if (dir && wsRef.current?.readyState === WebSocket.OPEN) {
        // Prevent sending duplicate directions
        if (dir === lastDirectionRef.current) {
          return;
        }
        
        lastDirectionRef.current = dir;
        const moveCommand = { type: "move", direction: dir };
        console.log("Sending move:", moveCommand);
        wsRef.current.send(JSON.stringify(moveCommand));
      }
    };

    window.addEventListener("keydown", handleKey);
    return () => window.removeEventListener("keydown", handleKey);
  }, []);

  // Send chat message
  const sendChat = useCallback(() => {
    if (!chatMessage.trim() || wsRef.current?.readyState !== WebSocket.OPEN) {
      return;
    }

    const chatData = {
      type: "chat",
      from: userId || "Unknown",
      message: chatMessage.trim(),
    };

    wsRef.current.send(JSON.stringify(chatData));
    setChatMessage("");
  }, [chatMessage, userId]);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-900 text-white p-4">
      <div className="max-w-4xl w-full">
        {/* Header */}
        <div className="mb-4 text-center">
          <h1 className="text-3xl font-bold mb-2">Snake Game</h1>
          <div className="flex items-center justify-center gap-4 text-sm">
            <span className={`px-3 py-1 rounded-full ${
              connectionStatus === "connected" 
                ? "bg-green-600" 
                : connectionStatus === "connecting"
                ? "bg-yellow-600"
                : "bg-red-600"
            }`}>
              {connectionStatus === "connected" && "Connected"}
              {connectionStatus === "connecting" && "Connecting..."}
              {connectionStatus === "disconnected" && "Disconnected"}
            </span>
            <span className="text-gray-400">Room: {gameId}</span>
            <span className="text-gray-400">Player: {userId}</span>
          </div>
          {error && (
            <div className="mt-2 text-red-400 text-sm">{error}</div>
          )}
        </div>

        {/* Game Canvas */}
        <div className="flex justify-center mb-4">
          <canvas
            ref={canvasRef}
            width={BOARD_WIDTH * CELL_SIZE}
            height={BOARD_HEIGHT * CELL_SIZE}
            className="border-4 border-gray-700 rounded shadow-2xl"
            style={{ background: "#111" }}
          />
        </div>

        {/* Game Info */}
        {gameState && (
          <div className="grid grid-cols-3 gap-4 mb-4 text-center">
            <div className="bg-gray-800 p-3 rounded">
              <div className="text-gray-400 text-sm">Score</div>
              <div className="text-2xl font-bold text-green-400">
                {gameState.playerSnake?.score?.value ?? 0}
              </div>
            </div>
            <div className="bg-gray-800 p-3 rounded">
              <div className="text-gray-400 text-sm">Length</div>
              <div className="text-2xl font-bold text-blue-400">
                {(gameState.playerSnake?.snakeBody?.length ?? 0) + 1}
              </div>
            </div>
            <div className="bg-gray-800 p-3 rounded">
              <div className="text-gray-400 text-sm">Opponents</div>
              <div className="text-2xl font-bold text-yellow-400">
                {gameState.otherSnakes?.length ?? 0}
              </div>
            </div>
          </div>
        )}

        {/* Controls Guide */}
        <div className="bg-gray-800 p-3 rounded mb-4 text-center text-sm">
          <span className="text-gray-400">Controls:</span>{" "}
          <span className="font-mono font-bold text-xl">↑ ↓ ← →</span> or{" "}
          <span className="font-mono font-bold text-xl">WASD</span> to move
        </div>

        {/* Chat */}
        <div className="bg-gray-800 rounded-lg p-4">
          <h2 className="text-lg font-semibold mb-2">💬 Chat</h2>
          <div className="h-40 overflow-y-auto bg-gray-900 p-3 rounded mb-2 border border-gray-700">
            {chatLog.length === 0 ? (
              <div className="text-gray-500 text-sm text-center mt-12">
                No messages yet
              </div>
            ) : (
              chatLog.map((msg, idx) => (
                <div
                  key={idx}
                  className={`mb-2 ${
                    msg.from === userId ? "text-green-400" : "text-blue-400"
                  }`}
                >
                  <span className="font-semibold">{msg.from}:</span>{" "}
                  <span className="text-white">{msg.message}</span>
                </div>
              ))
            )}
          </div>
          <div className="flex gap-2">
            <input
              className="flex-1 px-3 py-2 bg-gray-700 border border-gray-600 rounded text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={chatMessage}
              onChange={(e) => setChatMessage(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && sendChat()}
              placeholder="Type a message..."
              disabled={connectionStatus !== "connected"}
            />
            <button
              className="px-6 py-2 bg-blue-600 hover:bg-blue-700 rounded font-semibold transition disabled:bg-gray-600 disabled:cursor-not-allowed"
              onClick={sendChat}
              disabled={connectionStatus !== "connected" || !chatMessage.trim()}
            >
              Send
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SnakeGame;