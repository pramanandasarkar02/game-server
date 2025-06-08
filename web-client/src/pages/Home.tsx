
import React, { useState, useContext } from 'react';
import { useNavigate } from 'react-router-dom';
import { PlayerContext } from '../contexts/PlayerContext';
import type { Player } from '../types/player';
import { v4 as uuidv4 } from 'uuid';

const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:4000';

const Home: React.FC = () => {
  const { player, setPlayer } = useContext(PlayerContext);
  const [error, setError] = useState<string>('');
  const [response, setResponse] = useState<string>('');
  const [formData, setFormData] = useState({ name: player?.name || '', password: '' });
  const navigate = useNavigate();

  const createUserFormat = async (playerData: { name: string }) => {
    if (!playerData.name) {
      setError('Please enter a username');
      return;
    }

    const requestPayload: Player = {
      name: playerData.name,
      id: player?.id || uuidv4(),
      level: player?.level || 0,
    };

    try {
      const res = await fetch(`${baseURL}/connect`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(requestPayload),
      });
      const data = await res.json();
      if (res.ok) {
        setPlayer(requestPayload);
        setResponse(`Successfully connected to server: ${data.message}`);
        setError('');
        navigate('/games');
      } else {
        setError(`Failed to connect to server: ${data.message}`);
      }
    } catch (err) {
      setError(`Error connecting to server: ${err instanceof Error ? err.message : 'Unknown error'}`);
    }
  };

  const handleCreateUser = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    createUserFormat(formData);
  };

  const handleLogin = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    createUserFormat(formData);
  };

  return (
    <div className="min-h-screen flex justify-center items-center bg-gray-100">
      <div className="w-full max-w-md p-8 space-y-8 bg-white rounded-lg shadow-md">
        <div>
          <h2 className="text-2xl font-bold text-center text-gray-900">Create Account</h2>
          <div className="mt-8 space-y-6">
            <div className="rounded-md shadow-sm space-y-4">
              <div>
                <label htmlFor="username" className="block text-sm font-medium text-gray-700">
                  Username
                </label>
                <input
                  id="username"
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  placeholder="Enter username"
                />
              </div>
              <div>
                <label htmlFor="password" className="block text-sm font-medium text-gray-700">
                  Password
                </label>
                <input
                  id="password"
                  type="password"
                  value={formData.password}
                  onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                  className="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  placeholder="Enter password"
                />
              </div>
            </div>
            <button
              onClick={handleCreateUser}
              className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              Create Account
            </button>
          </div>
        </div>

        <div>
          <h2 className="text-2xl font-bold text-center text-gray-900">Login</h2>
          <div className="mt-8 space-y-6">
            <div className="rounded-md shadow-sm space-y-4">
              <div>
                <label htmlFor="login-username" className="block text-sm font-medium text-gray-700">
                  Username
                </label>
                <input
                  id="login-username"
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  placeholder="Enter username"
                />
              </div>
            </div>
            <button
              onClick={handleLogin}
              className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              Login
            </button>
          </div>
        </div>

        {error && <p className="text-red-500 text-sm text-center">{error}</p>}
        {response && <p className="text-green-500 text-sm text-center">{response}</p>}
      </div>
    </div>
  );
};

export default Home;
