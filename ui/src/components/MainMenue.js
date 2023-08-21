import { Button } from "antd";
import { useNavigate } from "react-router-dom";

function MainMenue() {
  const navigate = useNavigate();
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
      <Button
        className="ant-btn"
        type="primary"
        block
        onClick={() => navigate("/lobby")}
      >
        Lobby Suche
      </Button>
      <Button
        className="ant-btn"
        type="primary"
        block
        onClick={() => navigate("/kreis")}
      >
        Kreise
      </Button>
    </div>
  );
}

export default MainMenue;
