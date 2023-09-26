import { Button } from "antd";
import { useNavigate } from "react-router-dom";

function MainMenue(props) {
  const navigate = useNavigate();
  const isTokenUndefined = !props.token && props.token !== "" && props.token !== undefined
  return (
    <div className="main-menu-container">
      <Button
        className="ant-btn"
        type="primary"
        block
        onClick={() => alert("gibts noch nicht")}
      >
        Schnelles Spiel
      </Button>
      {!isTokenUndefined ? <Button
        className="ant-btn"
        type="primary"
        block
        onClick={() => navigate("/lobby")}
      >
        Lobby Suche
      </Button> : <></>}

      {!isTokenUndefined ? <Button
        className="ant-btn"
        type="primary"
        block
        onClick={() => navigate("/characterCreation")}
      >
        Charakter Erstellung
      </Button> : <></>}

      {isTokenUndefined ? <Button
        className="ant-btn"
        type="primary"
        block
        onClick={() => navigate("/login")} 
      >
        Einloggen
      </Button> : <></>}

      {!isTokenUndefined ? <Button
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
