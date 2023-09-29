import { Button } from "antd";
import { useNavigate } from "react-router-dom";

function MainMenue(props) {
  const navigate = useNavigate();
  return (
    <div className="main-menu-container">
      <Button
        className="ant-btn"
        type="primary"
        block
        onClick={() => navigate("/gameServer/random")}
      >
        Schnelles Spiel
      </Button>
      {!props.isTokenUndefined ? <Button
        className="ant-btn"
        type="primary"
        block
        onClick={() => navigate("/gameLobby")}
      >
        Lobby Suche
      </Button> : <></>}

      {!props.isTokenUndefined ? <Button
        className="ant-btn"
        type="primary"
        block
        onClick={() => navigate("/custom")}
      >
        Charaktererstellung
      </Button> : <></>}

      {props.isTokenUndefined ? <Button
        className="ant-btn"
        type="primary"
        block
        onClick={() => navigate("/login")} 
      >
        Einloggen
      </Button> : <></>}

      {!props.isTokenUndefined ? <Button
        className="ant-btn"
        type="primary"
        block
        onClick={() => props.removeToken()} 
      >
        Ausloggen
      </Button> : <></>}
    </div>
  );
}

export default MainMenue;
