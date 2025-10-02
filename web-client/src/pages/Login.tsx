import axios from "axios";
import { useContext, useState } from "react";
import PlayerContext from "../context/PlayerContext";
import type { Player } from "../types/player";

type LoginRequest = {
  username: string;
  password: string;
};

const Login = () => {
  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [message, setMessage] = useState<string>("");

  const {setPlayer } = useContext(PlayerContext);

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
    } catch (error: any) {
      console.error("Login error:", error);
      setMessage(error.response?.data?.message || "Login failed");
    }
  };

  return (
    <div>
      <form onSubmit={OnLoginButtonAction}>
        <div>
          <h1>Enter Username:</h1>
          <input
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="username"
          />
        </div>
        <div>
          <h1>Enter Password:</h1>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="password"
          />
        </div>
        <button type="submit">Login</button>
      </form>

      {/* Display login status */}
      {message && <p>{message}</p>}
      
    </div>
  );
};

export default Login;
