import React, { useState, useContext } from 'react';
import { useNavigate } from 'react-router-dom';
import { AuthContext } from '../contexts/AuthContext';


const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:4000';

interface RequestPayloadType {
  username: string;
  password: string;
}

const Home: React.FC = () => {
  const { player, setPlayer } = useContext(AuthContext);

  const [error, setError] = useState<string>('');
  const [response, setResponse] = useState<string>('');
  const [formData, setFormData] = useState({ username: player?.name || '', password: '' });
  const [isLogin, setIsLogin] = useState(true);
  const navigate = useNavigate();

  

  const connect = async (username: string, password: string) => {
    if (!username) {
      setError('Please enter a username');
      return;
    }
    if (!password) {
      setError('Please enter a password');
      return;
    }

    const requestPayload: RequestPayloadType = { username, password };

    try {
      const res = await fetch(`${baseURL}/connect`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(requestPayload),
      });
      const data = await res.json();
      console.log(data);
      if (res.ok) {
        setPlayer(data.player);
        
        setResponse(`Successfully logged in: ${data.message}`);
        setError('');
        navigate('/home');
      } else {
        setError(`Failed to log in: ${data.message}`);
      }
    } catch (err) {
      setError(`Error logging in: ${err instanceof Error ? err.message : 'Unknown error'}`);
    }
  };

  const register = async (username: string, password: string) => {
    if (!username) {
      setError('Please enter a username');
      return;
    }
    if (!password) {
      setError('Please enter a password');
      return;
    }

    const requestPayload: RequestPayloadType = { username, password };

    try {
      const res = await fetch(`${baseURL}/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(requestPayload),
      });
      const data = await res.json();
      console.log(data);
      if (res.ok) {
        setPlayer(data.player);
        // setToken(data.token);
        setResponse(`Successfully created account: ${data.message}`);
        setError('');
        navigate('/home');
      } else {
        setError(`Failed to create account: ${data.message}`);
      }
    } catch (err) {
      setError(`Error creating account: ${err instanceof Error ? err.message : 'Unknown error'}`);
    }
  };

  const handleLogin = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    connect(formData.username, formData.password);
  };

  const handleCreateAccount = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    register(formData.username, formData.password);
  };

  const toggleForm = () => {
    setIsLogin(!isLogin);
    setError('');
    setResponse('');
    setFormData({ username: player?.name || '', password: '' });
  };

  return (
    <div className="min-h-screen flex justify-center items-center bg-gray-100">
      <div className="w-full max-w-md p-8 space-y-8 bg-white rounded-lg shadow-md">
        <div>
          <h2 className="text-2xl font-bold text-center text-gray-900">
            {isLogin ? 'Login' : 'Create Account'}
          </h2>
          <div className="mt-8 space-y-6">
            <div className="rounded-md shadow-sm space-y-4">
              <div>
                <label htmlFor="username" className="block text-sm font-medium text-gray-700">
                  Username
                </label>
                <input
                  id="username"
                  type="text"
                  value={formData.username}
                  onChange={(e) => setFormData({ ...formData, username: e.target.value })}
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
              onClick={isLogin ? handleLogin : handleCreateAccount}
              className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              {isLogin ? 'Login' : 'Create Account'}
            </button>
            <button
              onClick={toggleForm}
              className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-indigo-600 bg-white hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              {isLogin ? 'Need an account? Create one' : 'Already have an account? Login'}
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