import React, { useState, useEffect, useContext, useRef } from 'react';
import { PlayerContext } from '../contexts/PlayerContext';
import type { Player, ChatMessage, TicTacToeState, WebSocketMessage } from '../types/player';

const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:4000';
const wsURL = import.meta.env.VITE_WS_URL || 'ws://localhost:4000';

const Games: React.FC = () => {
  const { player } = useContext(PlayerContext);
  const [matchID, setMatchID] = useState<string>('');
  const [gameState, setGameState] = useState<TicTacToeState | null>(null);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [chatInput, setChatInput] = useState<string>('');
  const [error, setError] = useState<string>('');
  const [status, setStatus] = useState<string>('Waiting to join queue...');
  const wsRef = useRef<WebSocket | null>(null);

  // Check server liveness
  const checkServerLiveness = async () => {
    try {
      const res = await fetch(`${baseURL}/ping`);
      const data = await res.json();
      if (res.ok) {
        setStatus(`Server is alive: ${data.message}`);
      } else {
        setError(`Server check failed: ${data.message}`);
      }
    } catch (err) {
      setError(`Error checking server: ${err instanceof Error ? err.message : 'Unknown error'}`);
    }
  };

  // Join queue
  const joinQueue = async () => {
    if (!player) return;
    try {
      const res = await fetch(`${baseURL}/queue/join`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ playerID: player.id, gameID: 'a2' }),
      });
      const data = await res.json();
      if (res.ok) {
        setStatus(`Joined queue: ${data.message}`);
        pollForMatch();
      } else {
        setError(`Failed to join queue: ${data.message}`);
      }
    } catch (err) {
      setError(`Error joining queue: ${err instanceof Error ? err.message : 'Unknown error'}`);
    }
  };

  // Poll for match
  const pollForMatch = async () => {
    if (!player) return;
    const interval = setInterval(async () => {
      try {
        const res = await fetch(`${baseURL}/match/${player.id}`);
        const data = await res.json();
        if (res.ok) {
          setMatchID(data.matchID);
          clearInterval(interval);
          joinGame(data.matchID);
        }
      } catch (err) {
        setError(`Error polling match: ${err instanceof Error ? err.message : 'Unknown error'}`);
      }
    }, 2000);
  };

  // Join game
  const joinGame = async (matchID: string) => {
    if (!player) return;
    try {
      const res = await fetch(`${baseURL}/running-match/${matchID}/${player.id}`);
      const data = await res.json();
      if (res.ok) {
        setMatchID(data.matchID);
        setStatus(`Joined game ${data.matchID} with players: ${data.players.join(', ')}`);
        startChat(data.matchID);
      } else {
        setError(`Failed to join game: ${data.message}`);
      }
    } catch (err) {
      setError(`Error joining game: ${err instanceof Error ? err.message : 'Unknown error'}`);
    }
  };

  // Handle Tic-Tac-Toe move
  const handleMove = (index: number) => {
    if (!player || !gameState || gameState.winner || gameState.isDraw || gameState.turn !== player.id || gameState.board[index]) {
      return;
    }
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      const move: WebSocketMessage = {
        type: 'move',
        data: { index, playerID: player.id },
      };
      wsRef.current.send(JSON.stringify(move));
    }
  };

  // Start chat and game WebSocket
  const startChat = (matchID: string) => {
    if (!player || !matchID) {
      setError('Cannot start chat: missing player or match ID');
      return;
    }
    wsRef.current = new WebSocket(`${wsURL}/chat/${matchID}/${player.id}`);
    wsRef.current.onopen = () => {
      setStatus(`Connected to game and chat for match ${matchID}`);
    };
    wsRef.current.onmessage = (event) => {
      const message: WebSocketMessage = JSON.parse(event.data);
      if (message.type === 'state') {
        setGameState(message.data as TicTacToeState);
        const state = message.data as TicTacToeState;
        if (state.winner) {
          setStatus(`Winner: ${state.winner}`);
        } else if (state.isDraw) {
          setStatus('Draw');
        } else {
          setStatus(`Turn: ${state.turn === player.id ? 'Your turn' : 'Opponent\'s turn'}`);
        }
      } else if (message.type === 'chat') {
        setMessages((prev) => [...prev, message.data as ChatMessage]);
      }
    };
    wsRef.current.onerror = (err) => {
      setError(`WebSocket error: ${err}`);
    };
    wsRef.current.onclose = () => {
      setStatus('Disconnected from game and chat');
    };
    // Periodic ping
    const pingInterval = setInterval(() => {
      if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify({ type: 'ping' }));
      }
    }, 30000);
    return () => clearInterval(pingInterval);
  };

  // Send chat message
  const sendMessage = () => {
    if (!chatInput.trim() || !wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) return;
    const message: WebSocketMessage = {
      type: 'chat',
      data: {
        matchId: matchID,
        playerId: player!.id,
        content: chatInput,
        timestamp: new Date().toISOString(),
      },
    };
    wsRef.current.send(JSON.stringify(message));
    setChatInput('');
  };

  // Cleanup WebSocket on unmount
  useEffect(() => {
    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, []);

  if (!player) {
    return <div className="text-center text-red-500">No Player Data Available</div>;
  }

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col items-center p-4">
      <div className="w-full max-w-4xl bg-white rounded-lg shadow-md p-6 space-y-6">
        <div className="flex gap-4">
          <h1 className="text-2xl font-bold">Level: {player.level}</h1>
          <h2 className="text-xl font-semibold">Name: {player.name}</h2>
          <h4 className="text-lg text-gray-600">ID: {player.id}</h4>
        </div>

        <div className="flex gap-4">
          <button
            onClick={checkServerLiveness}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            Check Server
          </button>
          <button
            onClick={joinQueue}
            className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
          >
            Join Queue
          </button>
        </div>

        {gameState && (
          <div className="grid grid-cols-3 gap-2 w-64 mx-auto">
            {gameState.board.map((cell, index) => (
              <button
                key={index}
                onClick={() => handleMove(index)}
                className="h-16 bg-gray-200 text-2xl font-bold flex items-center justify-center hover:bg-gray-300 disabled:bg-gray-400"
                disabled={
                  Boolean(cell) || 
                  Boolean(gameState.winner) || 
                  Boolean(gameState.isDraw) || 
                  (gameState.turn !== player.id)
                }
              >
                {cell}
              </button>
            ))}
          </div>
        )}
        <p className="text-center text-lg">{status}</p>

        <div className="mt-8">
          <h3 className="text-lg font-medium">Chat</h3>
          <div className="h-64 overflow-y-auto border border-gray-300 p-2">
            {messages.map((msg, index) => (
              <p key={index} className={msg.playerId === player.id ? 'text-right text-blue-600' : 'text-left'}>
                [{new Date(msg.timestamp).toLocaleString()}] {msg.playerId}: {msg.content}
              </p>
            ))}
          </div>
          <div className="flex gap-2 mt-2">
            <input
              type="text"
              value={chatInput}
              onChange={(e) => setChatInput(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
              className="flex-1 px-3 py-2 border border-gray-300 rounded"
              placeholder="Type a message..."
            />
            <button
              onClick={sendMessage}
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Send
            </button>
          </div>
        </div>

        {error && <p className="text-red-500 text-sm text-center">{error}</p>}
      </div>
    </div>
  );
};

export default Games;