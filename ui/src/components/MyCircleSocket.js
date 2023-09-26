import { useEffect, useRef, useState } from "react";
import { useWebSocket } from "react-use-websocket/dist/lib/use-websocket";
import useWindowDimensions from "../hooks/useWindowDimensions";

let userPlayer = '';

const colors = ["red", "green", "blue", "yellow", "maroon", "purple", "lime", "olive", "teal", "aqua"]

function createCircle(Pcx, Pcy, Pcolor, name, size) {
  return <>
    <circle id="1" key={"c"+name} cx={Pcx} cy={Pcy} r={size} stroke="black" stroke-width="3" fill={colors[Pcolor]} />
    <text key={"t"+name} x={Pcx - 50} y={Pcy + 20 + size} fontSize="20" fill="black">{name}</text>
  </>
}

function createCircleNpc(Pcx, Pcy, Pcolor) {
  return <circle key={"c"+Pcx+"-"+Pcy} cx={Pcx} cy={Pcy} r="10" stroke={colors[Pcolor]} stroke-width="3" fill={colors[Pcolor]} />
}

function MyCircleSocket(props) {
  const [cameraPosition, setCameraPosition] = useState({ x: 0, y: 0 });
  const svgRef = useRef()
  const { sendMessage, lastMessage } = useWebSocket(props.serverUrl);
  const [playerObject, setPlayerObject] = useState([]);
  const [otherPlayerObjects, setOtherPlayerObjects] = useState([]);
  const [npcObjects, setNpcObjects] = useState([]);
  const {width, height} = useWindowDimensions();

  useEffect(() => {
    if (lastMessage !== null) {
      const messageData = JSON.parse(lastMessage.data)
      userPlayer = messageData.player     //folgt zurzeit dem ersten Spieler: später nach ID   //.find(player => player.id === ''); // Replace userId with the actual user's ID
    if (userPlayer) {
      // Adjust the camera position to follow the user player
      setCameraPosition({
        x: userPlayer.x - width / 2, // Adjust as needed
        y: userPlayer.y - height / 2, // Adjust as needed
      });
    }
      setPlayerObject(messageData.player)
      setOtherPlayerObjects(messageData.otherPlayer)
      setNpcObjects(messageData.npc)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [lastMessage]);

  const fullScreen = {
    position: "fixed",
    top: 0,
    bottom: 0,
    left: 0,
    right: 0,
    // backgroundColor: 'black'
  };

  useEffect(() => {
    setPositionUpdater()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  function setPositionUpdater() {
    if (svgRef.current) {
      
      
      const circleCenterX = width/2;
      const circleCenterY = height/2;
  
      let relativeMousePosition = { x: 0, y: 0 };
  
      document.onmousemove = (e) => {
        relativeMousePosition.x = e.clientX - circleCenterX;
        relativeMousePosition.y = e.clientY - circleCenterY;
      };
  
      document.ontouchmove = (e) => {
        relativeMousePosition.x = e.touches[0].clientX - circleCenterX;
        relativeMousePosition.y = e.touches[0].clientY - circleCenterY;
      };
  
      setInterval(() => updatePosition(relativeMousePosition), 10);
    
  }
  }

  function updatePosition(mousePosition) {
    // Calculate the player's new position based on the mouse/touch position
    const newX = mousePosition.x;
    const newY = mousePosition.y;
  
    // Send the updated position to the server
    sendMessage(JSON.stringify({ x: newX, y: newY }));
  }

  return (
    <div>
      <svg className="karo-container" bac ref={svgRef} style={fullScreen} width={width} height={height}>
        <rect x={-10 - cameraPosition.x} y={-10 - cameraPosition.y} fill="white" width={10} height={5010} />
        <rect x={5000 - cameraPosition.x} y={-10 - cameraPosition.y} fill="white" width={10} height={5010} />
        <rect x={0 - cameraPosition.x} y={-10 - cameraPosition.y} fill="white" width={5010} height={10} />
        <rect x={-10 - cameraPosition.x} y={5000 - cameraPosition.y} fill="white" width={5020} height={10} />
        {npcObjects?.map(obj =>
           createCircleNpc(obj.x - cameraPosition.x, obj.y - cameraPosition.y, obj.color)
        )}
        {otherPlayerObjects?.map(obj =>
          createCircle(obj.x - cameraPosition.x, obj.y - cameraPosition.y, obj.color, obj.id, obj.size)
          )}
        {createCircle(playerObject.x - cameraPosition.x, playerObject.y - cameraPosition.y, playerObject.color, playerObject.id, playerObject.size)}
      </svg>
    </div>
  );
}

export default MyCircleSocket;
