import { Route, Routes, Link, useLocation } from "react-router-dom";
import MainMenue from "./components/MainMenue";
import Lobby from "./components/Lobby";
import MyCircle from "./components/MyCircle";
import GameContainer from "./components/GameContainer";
import CharacterCreation from "./components/CharacterCreation";
import Login from "./components/Login";
import { useEffect, useState } from "react"; // Hook hinzufügen

function App() {
  const location = useLocation();
  const [hideHeaderAndButton, setHideHeaderAndButton] = useState(false);

  useEffect(() => {
    // Überprüfen, ob sich der Benutzer auf der Seite "/game" befindet
    if (location.pathname.startsWith("/game")) {
      setHideHeaderAndButton(true);
    } else {
      setHideHeaderAndButton(false);
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
      <Routes>
        <Route path="/" element={<MainMenue />} />
        <Route path="/lobby" element={<Lobby/>}/>
        <Route path="/game/:id" element={<GameContainer />}/>
        <Route path="/kreis" element={<MyCircle />}/>
        <Route path="/CharacterCreation" element={<CharacterCreation />}/>
        <Route path="/login" element={<Login />} />
      </Routes>
    </div>
  );
}

export default App;
