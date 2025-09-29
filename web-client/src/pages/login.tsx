import axios from "axios";
import { useState } from "react";

type LoginRequest = {
  username: string;
  password: string;
};

const Login = () => {
  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [message, setMessage] = useState<string>("");

  const OnLoginButtonAction = async (e: React.FormEvent) => {
    e.preventDefault();

    const loginData: LoginRequest = { username, password };

    try {
      const response = await axios.post("http://localhost:8080/api/login", loginData);
      const data = response.data;
      setMessage(data.message || "Login successful");
    } catch (error: any) {
      console.error("Login error:", error);
      setMessage("Failed to connect to the server");
    }
  };

  return (
    <div>
      <form onSubmit={OnLoginButtonAction}>
        <div>
          <h1>Enter Username: </h1>
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
      {message && <p>{message}</p>}
    </div>
  );
};

export default Login;
