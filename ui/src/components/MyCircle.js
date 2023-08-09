import { useEffect, useRef, useState } from "react";
import { normalize } from "../helper/PositionHelper";

const width = window.innerWidth;
const height = window.innerHeight;

function createCircle(Pcx, Pcy, Pcolor, name) {
  return <>
    <circle cx={Pcx} cy={Pcy} r="40" stroke="white" stroke-width="3" fill={Pcolor} />
    <text x={Pcx - 50} y={Pcy + 70} fontSize="20" fill="white">{name}</text>
  </>
}

function MyCircle() {
  const svgRef = useRef()
  const [playerPos, setPlayerPos] = useState({ x: 100, y: 100 });



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
      setInterval(() => updatePosition(point), 20);
    }
  }

  function updatePosition(point) {
    const svgElement = svgRef.current
    if(svgElement) {
      const screenCTM = svgElement.getScreenCTM();
      if(screenCTM) {
        const loc = point.matrixTransform(screenCTM.inverse());
        const normalized = normalize(loc.x - width / 2, loc.y - height / 2);

        setPlayerPos(prevPlayerPos => {
          const newX = prevPlayerPos.x + normalized.x;
          const newY = prevPlayerPos.y + normalized.y;
          return { x: newX, y: newY };
      });
      }
    }
  }

  return (
    <div>
      <svg ref={svgRef} style={fullScreen} width={width} height={height}>
        {createCircle(playerPos.x, playerPos.y, "fuchsia", "player")}
        {createCircle(500, 700, "red", "bot")}
      </svg>
    </div>
  );
}

export default MyCircle;
