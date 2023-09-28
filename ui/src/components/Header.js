import { Link, useLocation } from "react-router-dom";
import HighscoreList from './HighscoreList';
import { useEffect, useState } from "react";
import { doGetRequest } from "../helper/RequestHelper";

function Header() {  
    const location = useLocation();
    const [hideHeaderAndButton, setHideHeaderAndButton] = useState(false);
    const [hideHighscoreList, setHideHighscoreList] = useState(false);
    const [highscores, setHighscores] = useState(null);
  
    useEffect(() => {
      if (location.pathname.startsWith("/gameServer")) {
        setHideHeaderAndButton(true);
        setHideHighscoreList(true);
      } else {
        setHideHeaderAndButton(false);
        setHideHighscoreList(false);
        doGetRequest('highscore').then(
          res => {
            setHighscores(
              res.data.map(row => ({
                highscore: row.Highscore,
                name: row.Name
              }))
            );
          }
        )
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
          <HighscoreList highscores={highscores} />
        </div>
      )}
  </div>);
}

export default Header;
