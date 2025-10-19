import axios from "axios";
import { useContext, useState } from "react";
import PlayerContext from "../context/PlayerContext";
import type { Player } from "../types/player";
import { useNavigate } from "react-router-dom";

type LoginRequest = {
  username: string;
  password: string;
};

const Login = () => {
  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [message, setMessage] = useState<string>("");
  const navigate = useNavigate();
  const { setPlayer } = useContext(PlayerContext);

  const OnLoginButtonAction = async (e: React.FormEvent) => {
    e.preventDefault();

    const loginData: LoginRequest = { username, password };

    try {
      const response = await axios.post("http://localhost:8080/api/login", loginData);
      const data = response.data;

      const newPlayer: Player = {
        username: data.username,
        userId: data.userId,
        playerStatus: data.playerStatus,
      };
      setPlayer(newPlayer);

      setMessage(data.message || "Login successful");
      navigate("/");
    } catch (error: any) {
      console.error("Login error:", error);
      setMessage(error.response?.data?.message || "Login failed");
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900">
      <div className="bg-gray-800 rounded-2xl shadow-2xl p-8 w-full max-w-md text-white">
        <h1 className="text-3xl font-bold text-center mb-6">Login</h1>
        <form onSubmit={OnLoginButtonAction} className="space-y-6">
          <div>
            <label className="block text-sm font-medium mb-2">Username</label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Enter your username"
              className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Enter your password"
              className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <button
            type="submit"
            className="w-full py-2 bg-blue-600 hover:bg-blue-700 rounded-xl text-lg font-semibold transition-all"
          >
            Login
          </button>
        </form>

        <div className="text-center mt-6">
          <p className="text-gray-400 text-sm">
            Don’t have an account?{" "}
            <span
              onClick={() => navigate("/signup")}
              className="text-blue-400 hover:text-blue-300 cursor-pointer"
            >
              Sign Up
            </span>
          </p>
        </div>

        {message && (
          <p
            className={`mt-4 text-center text-sm ${
              message.includes("successful") ? "text-green-400" : "text-red-400"
            }`}
          >
            {message}
          </p>
        )}
      </div>
    </div>
  );
};

export default Login;
