import React, { useEffect, useRef } from 'react';

function createCircle(Pcx, Pcy, Pcolor, name, size) {
  return <>
    <circle id="1" key={"c"+name} cx={Pcx} cy={Pcy} r={size} stroke="black" strokeWidth="3" fill={Pcolor} />
    <text key={"name"+name} x={Pcx} y={Pcy + 20 + size} fontSize="20" fill="black" textAnchor="middle">{name}</text>
    <text key={"size"+name} x={Pcx} y={Pcy+5} fontSize="20" fill="black" textAnchor="middle">{Math.round(size)}</text>
  </>
}

function createCircleNpc(Pcx, Pcy, Pcolor, idx) {
  return <circle key={"npc"+idx} cx={Pcx} cy={Pcy} r="10" stroke={Pcolor} strokeWidth="3" fill={"#191a17"} />
}

function Gamefield(props) {
    const svgRef = useRef()

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
        
        
        const circleCenterX = props.width/2;
        const circleCenterY = props.height/2;
    
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
      props.sendMessage(JSON.stringify({ mode: "pos", x: newX, y: newY }));
      // sendMessage(JSON.stringify({ mode: "msg", message: newX + " " + newY }));
    }

  return (
    <svg className="karo-container" ref={svgRef} style={fullScreen} width={props.width} height={props.height}>
      {props.playerObject.length !== 0 ? <>
        <rect x={-10 - props.cameraPosition.x} y={-10 - props.cameraPosition.y} fill="white" width={10} height={5010} />
        <rect x={5000 - props.cameraPosition.x} y={-10 - props.cameraPosition.y} fill="white" width={10} height={5010} />
        <rect x={0 - props.cameraPosition.x} y={-10 - props.cameraPosition.y} fill="white" width={5010} height={10} />
        <rect x={-10 - props.cameraPosition.x} y={5000 - props.cameraPosition.y} fill="white" width={5020} height={10} />
        {props.npcObjects?.map((obj, idx) =>
           createCircleNpc(obj.x - props.cameraPosition.x, obj.y - props.cameraPosition.y, obj.color, idx)
        )}
        {props.otherPlayerObjects?.map(obj =>
          createCircle(obj.x - props.cameraPosition.x, obj.y - props.cameraPosition.y, obj.color, obj.name, obj.size)
          )}
        {createCircle(props.playerObject.x - props.cameraPosition.x, props.playerObject.y - props.cameraPosition.y, props.playerObject.color, props.playerObject.name, props.playerObject.size)}
        </>: <></>}
      </svg>
  );

}

export default Gamefield;
