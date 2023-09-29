import { useEffect, useState } from "react";
import { useWebSocket } from "react-use-websocket/dist/lib/use-websocket";
import useWindowDimensions from "../hooks/useWindowDimensions";
import HighscoreList from "./HighscoreList";
import Chat from "./Chat";
import Gamefield from "./Gamefield";

function WebsocketContainer(props) {
  const [cameraPosition, setCameraPosition] = useState({ x: 0, y: 0 });
  const { sendMessage, lastMessage } = useWebSocket(props.serverUrl + "/" + (props.token ? props.token : "undefined"));
  const [playerObject, setPlayerObject] = useState([]);
  const [otherPlayerObjects, setOtherPlayerObjects] = useState([]);
  const [npcObjects, setNpcObjects] = useState([]);
  const [chatMessages, setChatMessages] = useState([]);
  const [gameServerScores, setGameServerScores] = useState([]);
  const { width, height } = useWindowDimensions();

  useEffect(() => {
    if (lastMessage !== null) {
      const messageData = JSON.parse(lastMessage.data)
      const userPlayer = messageData.player
      if (userPlayer) {
        // Adjust the camera position to follow the user player
        setCameraPosition({
          x: userPlayer.x - width / 2,
          y: userPlayer.y - height / 2,
        });
      }
      setPlayerObject(messageData.player)
      setOtherPlayerObjects(messageData.otherPlayer)
      if (messageData.npc) {
        setNpcObjects(messageData.npc)
      }
      if (messageData.score) {
        setGameServerScores(messageData.score)
      }

      if (messageData.message) {
        for (let msg of messageData.message) {
          setChatMessages([...chatMessages, msg]);
        }
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [lastMessage]);

  return (<>
    <Gamefield cameraPosition={cameraPosition} playerObject={playerObject} otherPlayerObjects={otherPlayerObjects} npcObjects={npcObjects} width={width} height={height} sendMessage={sendMessage} />
    <div className="highscore-list">
      👑 Gameserver Scores 👑 <br />
      • Name, Score: <br />
      <HighscoreList highscores={gameServerScores} />
    </div>
    <Chat chatMessages={chatMessages} sendMessage={sendMessage} />
  </>
  );
}

export default WebsocketContainer;
