import { useEffect, useRef, useState } from "react";
import { useWebSocket } from "react-use-websocket/dist/lib/use-websocket";

const width = window.innerWidth;
const height = window.innerHeight;

const colors = ["red", "green", "blue", "yellow", "silver", "maroon", "purple", "lime", "olive", "teal", "aqua"]

function createCircle(Pcx, Pcy, Pcolor, name) {
  return <>
    <circle key={"c"+name} cx={Pcx} cy={Pcy} r="40" stroke="white" stroke-width="3" fill={colors[Pcolor]} />
    <text key={"t"+name} x={Pcx - 50} y={Pcy + 70} fontSize="20" fill="white">{name}</text>
  </>
}

function MyCircleSocket(props) {
  const svgRef = useRef()
  const { sendMessage, lastMessage, readyState } = useWebSocket(props.serverUrl);
  const [playerPos, setPlayerPos] = useState({ x: 100, y: 100 });
  const [gameObjects, setGameObjects] = useState([]);

  useEffect(() => {
    if (lastMessage !== null) {
      const messageData = JSON.parse(lastMessage.data)
      const e = Object.keys(messageData).map((e, i) => {
        return {key: e, x: messageData[e].x, y: messageData[e].y, color: messageData[e].color}
      })
      setGameObjects(e)
    }
  }, [lastMessage]);

  const fullScreen = {
    position: "fixed",
    top: 0,
    bottom: 0,
    left: 0,
    right: 0,
    backgroundColor: 'black'
  };

  useEffect(() => {
    setPositionUpdater()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  function setPositionUpdater() {
    if (svgRef.current) {
      let point = svgRef.current.createSVGPoint();
      document.onmousemove = (e) => {
        point.x = e.clientX;
        point.y = e.clientY;
      };
      document.ontouchmove = (e) => {
        point.x = e.touches[0].clientX;
        point.y = e.touches[0].clientY;
      };
      setInterval(() => updatePosition(point), 10);
    }
  }

  function updatePosition(point) {
    setPlayerPos({x: point.x, y: point.y})    
    sendMessage(JSON.stringify({x: point.x, y: point.y}))
  }

  return (
    <div>
      <svg ref={svgRef} style={fullScreen} width={width} height={height}>
        {/* {createCircle(playerPos.x, playerPos.y, "blue", "player")}
        {createCircle(500, 700, "red", "bot")} */}
        {gameObjects.map(obj =>
          createCircle(obj.x, obj.y, obj.color, obj.key)
          )}
      </svg>
    </div>
  );
}

export default MyCircleSocket;
