import { Link, useLocation } from "react-router-dom";
import HighscoreList from './HighscoreList';
import { useEffect, useState } from "react";

function Header() {  
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
  
  return (<div>
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
  </div>);
}

export default Header;
