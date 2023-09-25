import React, { useEffect, useState } from 'react';
import { doGetRequest } from "../helper/RequestHelper";

function HighscoreList() {
  const [highscores, setHighscores] = useState(null); // Null als Initialwert


  useEffect(() => {
    doGetRequest('highscores').then(
      res => {
        setHighscores(
          res.data.map(row => ({
            highscore: row.Highscore,
            name: row.Name
          }))
        );
      }
    )
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Wenn die Daten noch geladen werden, zeige eine Ladeanzeige
  if (highscores === null) {
    return <div>Loading... </div>;
  }

  // Zeige die gesamte JSON-Antwort
  return (
    <ul className="highscore-list">
      {highscores.map((score, index) => (
        <li key={index}>
          {score.name}, {score.highscore}
        </li>
      ))}
    </ul>
  );
}

export default HighscoreList;




