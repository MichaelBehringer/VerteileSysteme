import { Route, Routes } from "react-router-dom";
import MainMenue from "./components/MainMenue";
import Lobby from "./components/Lobby";
import GameContainer from "./components/GameContainer";
import CharacterCreation from "./components/CharacterCreation";
import Login from "./components/Login";
import Header from "./components/Header";

function App() {
  return (
    <div className="app-container">
      <Header />
      <Routes>
        <Route path="/" element={<MainMenue />} />
        <Route path="/lobby" element={<Lobby/>}/>
        <Route path="/game/:id" element={<GameContainer />}/>
        <Route path="/CharacterCreation" element={<CharacterCreation />}/>
        <Route path="/login" element={<Login />} />
      </Routes>
    </div>
  );
}

export default App;
