import React, { useState } from "react";
import { Button, Input } from "antd";
import { useNavigate } from "react-router-dom";

function Login() {
  const navigate = useNavigate();

  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  // Funktion, die aufgerufen wird, wenn der Benutzer sich einloggt
  const handleLogin = () => {
    // Hier kannst du die Logik für den Login-Vorgang implementieren
    // Zum Beispiel: Überprüfung der Anmeldeinformationen, Zustand des Benutzer-Logins usw.
    // Für dieses Beispiel setzen wir isLoggedIn auf true.
    setIsLoggedIn(true);
  };

  if (isLoggedIn) {
    navigate("/"); // Leite den Benutzer zur Hauptseite (MainMenue) zurück, wenn eingeloggt
  }

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
