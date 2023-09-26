import React, { useState } from "react";
import { Button, Input } from "antd";
import { useNavigate } from "react-router-dom";
import { doCustomPostRequest } from "../helper/RequestHelper";
import { myToastError } from "../helper/ToastHelper";

function Login(props) {
  const navigate = useNavigate();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  // Funktion, die aufgerufen wird, wenn der Benutzer sich einloggt
  const handleLogin = () => {
    const params = {username: username, password: password};
		doCustomPostRequest("auth/token", params).then((response) => {
			props.setToken(response.data.accessToken);
      navigate("/")
		}, error => {
			if (error.response.status === 401) {
				myToastError("Benutzername oder Passwort falsch!");
			}
			return error;
		});
  };

  return (
    <div className="login-container">
      <h2 style={{ color: "white" }}>Login</h2>
      <Input
        placeholder="Benutzername"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        style={{ marginBottom: "10px" }}
      />
      <Input
        type="password"
        placeholder="Passwort"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        style={{ marginBottom: "10px" }}
      />
      <Button type="primary" onClick={handleLogin}>
        Einloggen
      </Button>
    </div>
  );
}

export default Login;
