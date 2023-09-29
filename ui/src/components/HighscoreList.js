import React from 'react';

function HighscoreList(props) {
  if (props.highscores === null) {
    return <div>Loading... </div>;
  }

  // Zeige die gesamte JSON-Antwort
  return (
    <ul className="highscore-list">
      {props.highscores.map((score, index) => (
        <li className='chat-msg' key={index}>
          {score.name}, {Math.round(score.highscore)}
        </li>
      ))}
    </ul>
  );
}

export default React.memo(HighscoreList);
