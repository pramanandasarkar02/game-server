import axios from "axios";
import { useContext, useState } from "react";
import PlayerContext from "../context/PlayerContext";
import type { Player } from "../types/player";
import { useNavigate } from "react-router-dom";

type SignupRequest = {
  username: string;
  password: string;
};

const Signup = () => {
  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [message, setMessage] = useState<string>("");
  const navigate = useNavigate();
  const { setPlayer } = useContext(PlayerContext);

  const OnSignupButtonAction = async (e: React.FormEvent) => {
    e.preventDefault();

    const signupData: SignupRequest = { username, password };

    try {
      const response = await axios.post("http://localhost:8080/api/signup", signupData);
      const data = response.data;

      const newPlayer: Player = {
        username: data.username,
        userId: data.userId,
        playerStatus: data.playerStatus,
      };
      setPlayer(newPlayer);

      setMessage(data.message || "Signup successful");
      navigate("/");
    } catch (error: any) {
      console.error("Signup error:", error);
      setMessage(error.response?.data?.message || "Signup failed");
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900">
      <div className="bg-gray-800 rounded-2xl shadow-2xl p-8 w-full max-w-md text-white">
        <h1 className="text-3xl font-bold text-center mb-6">Sign Up</h1>

        <form onSubmit={OnSignupButtonAction} className="space-y-6">
          <div>
            <label className="block text-sm font-medium mb-2">Username</label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Enter your username"
              className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-xl focus:outline-none focus:ring-2 focus:ring-green-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Enter your password"
              className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-xl focus:outline-none focus:ring-2 focus:ring-green-500"
            />
          </div>

          <button
            type="submit"
            className="w-full py-2 bg-green-600 hover:bg-green-700 rounded-xl text-lg font-semibold transition-all"
          >
            Sign Up
          </button>
        </form>

        <div className="text-center mt-6">
          <p className="text-gray-400 text-sm">
            Already have an account?{" "}
            <span
              onClick={() => navigate("/login")}
              className="text-green-400 hover:text-green-300 cursor-pointer"
            >
              Login
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

export default Signup;
