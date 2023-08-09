import { useEffect, useRef, useState } from "react";
import { useWebSocket } from "react-use-websocket/dist/lib/use-websocket";

const width = window.innerWidth;
const height = window.innerHeight;

const colors = ["red", "green", "blue", "yellow", "silver", "maroon", "purple", "lime", "olive", "teal", "aqua"]

function createCircle(Pcx, Pcy, Pcolor, name, size) {
  return <>
    <circle key={"c"+name} cx={Pcx} cy={Pcy} r={size} stroke="white" stroke-width="3" fill={colors[Pcolor]} />
    <text key={"t"+name} x={Pcx - 50} y={Pcy + 20 + size} fontSize="20" fill="white">{name}</text>
  </>
}

function createCircleNpc(Pcx, Pcy, Pcolor, name) {
  return <circle key={"c"+name} cx={Pcx} cy={Pcy} r="10" stroke={colors[Pcolor]} stroke-width="3" fill={colors[Pcolor]} />
}

function MyCircleSocket(props) {
  const svgRef = useRef()
  const { sendMessage, lastMessage, readyState } = useWebSocket(props.serverUrl);
  const [playerPos, setPlayerPos] = useState({ x: 100, y: 100 });
  const [playerObjects, setPlayerObjects] = useState([]);
  const [npcObjects, setNpcObjects] = useState([]);

  useEffect(() => {
    if (lastMessage !== null) {
      const messageData = JSON.parse(lastMessage.data)
      const objs = Object.keys(messageData)
      const e = objs.slice(1).map((e, i) => {
        return {key: e, x: messageData[e].x, y: messageData[e].y, color: messageData[e].c, size: messageData[e].s}
      })
      if(objs[0]==="00000000-0000-0000-0000-000000000000") {
        setPlayerObjects(e)
      } else {
        setNpcObjects(e)
      }
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
        {npcObjects.map(obj =>
          createCircleNpc(obj.x, obj.y, obj.color, obj.key)
          )}
          {playerObjects.map(obj =>
            createCircle(obj.x, obj.y, obj.color, obj.key, obj.size)
            )}
      </svg>
    </div>
  );
}

export default MyCircleSocket;
