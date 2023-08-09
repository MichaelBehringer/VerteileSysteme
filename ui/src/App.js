import {Route, Routes} from "react-router-dom";
import MainMenue from "./components/MainMenue";
import Lobby from "./components/Lobby";
import MyCircle from "./components/MyCircle";
import GameContainer from "./components/GameContainer";

function App() {
  return (    
    <div>
      <h1>Hallo zu unserem spiel 🦚</h1>
      <Routes>
        <Route path="/" element={<MainMenue />} />
        <Route path="/lobby" element={<Lobby/>}/>
        <Route path="/game/:id" element={<GameContainer />}/>
        <Route path="/kreis" element={<MyCircle />}/>
      </Routes>
      </div>
  );
}

export default App;
