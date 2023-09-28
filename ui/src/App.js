import { Route, Routes } from "react-router-dom";
import MainMenue from "./components/MainMenue";
import Lobby from "./components/Lobby";
import GameContainer from "./components/GameContainer";
import CharacterCreation from "./components/CharacterCreation";
import Login from "./components/Login";
import Header from "./components/Header";
import useToken from "./hooks/useToken";

function App() {
	const {token, removeToken, setToken} = useToken();
  const isTokenUndefined = !token && token !== "" && token !== undefined
  return (
    <div className="app-container">
      <Header />
      <Routes>
        <Route path="/" element={<MainMenue isTokenUndefined={isTokenUndefined} removeToken={removeToken}/>} />
        {!isTokenUndefined ? <Route path="/gameLobby" element={<Lobby/>}/> : <></>}
        <Route path="/gameServer/:id" element={<GameContainer token={token}/>}/>
        {!isTokenUndefined ? <Route path="/custom" element={<CharacterCreation token={token}/>}/> : <></>}
        <Route path="/login" element={<Login setToken={setToken} />} />
      </Routes>
    </div>
  );
}

export default App;
