import React, { useState } from "react";
import { Button, Input, Modal } from "antd";
import { useNavigate } from "react-router-dom";
import { doCustomPostRequest, doPutRequestAuth } from "../helper/RequestHelper";
import { myToastError, myToastSuccess } from "../helper/ToastHelper";

function Login(props) {
  const navigate = useNavigate();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  
  const [usernameReg, setUsernameReg] = useState("");
  const [passwordReg, setPasswordReg] = useState("");
  const [isModalOpen, setIsModalOpen] = useState(false);

  const handleReg = () => {
    const params = {username: usernameReg, password: passwordReg};
    doPutRequestAuth("user", params, props.token).then(() => {
      myToastSuccess('Anlegen erfolgreich');
      setUsernameReg("")
      setPasswordReg("")
      setIsModalOpen(false);
		}, error => {
			if (error.response.status === 400) {
				myToastError("Username bereits vorhanden!");
			} else {
        myToastError("Fehler!");
      }
			return error;
		});
  };

  const openModal = () => {
    setUsername("")
    setPassword("")
    setIsModalOpen(true);
  };

  const handleLogin = () => {
    const params = {username: username, password: password};
		doCustomPostRequest("auth/token", params).then((response) => {
			props.setToken(response.data.accessToken);
      myToastSuccess('Hallo ' + username + '!');
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
      <Modal title="Registrieren" open={isModalOpen} onOk={()=>handleReg()} onCancel={()=>setIsModalOpen(false)}>
      <Input
        placeholder="Benutzername"
        value={usernameReg}
        onChange={(e) => setUsernameReg(e.target.value)}
        style={{ marginBottom: "10px" }}
      />
      <Input
        type="password"
        placeholder="Passwort"
        value={passwordReg}
        onChange={(e) => setPasswordReg(e.target.value)}
        style={{ marginBottom: "10px" }}
      />
      </Modal>
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
      <Button type="primary" onClick={()=>openModal()}>
        Registrieren
      </Button>
    </div>
  );
}

export default Login;
