import { Route, Routes, Link, useLocation } from "react-router-dom";
import MainMenue from "./components/MainMenue";
import Lobby from "./components/Lobby";
import GameContainer from "./components/GameContainer";
import CharacterCreation from "./components/CharacterCreation";
import Login from "./components/Login";
import HighscoreList from './components/HighscoreList';
import { useEffect, useState } from "react";

function App() {
  const location = useLocation();
  const [hideHeaderAndButton, setHideHeaderAndButton] = useState(false);
  const [hideHighscoreList, setHideHighscoreList] = useState(false);

  useEffect(() => {
    if (location.pathname.startsWith("/game")) {
      setHideHeaderAndButton(true);
      setHideHighscoreList(true);
    } else {
      setHideHeaderAndButton(false);
      setHideHighscoreList(false);
    }
  }, [location.pathname]);

  return (
    <div className="app-container">
      {!hideHeaderAndButton && (
        <>
          <Link to="/" className="back-to-main-button">
            🦚
          </Link>
          <h1>Hallo zu unserem Spiel 🦚</h1>
        </>
      )}
      {!hideHighscoreList && (
        <div className="highscore-list">
          👑 Top Highscores 👑 <br/>
          • Name, Highscore:  <br/>
          <HighscoreList  />
        </div>
      )}
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
