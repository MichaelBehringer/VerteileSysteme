import { Button } from "antd";
import { useNavigate } from "react-router-dom";

function MainMenue() {
    const navigate = useNavigate();
  return (    
    <div>
        <Button style={{margin: '10px'}} type="primary" block onClick={()=>alert('gibts noch nicht')}>Schnelles spiel</Button>
        <Button style={{margin: '10px'}} type="primary" block onClick={()=>navigate("/lobby")}>Lobby suche</Button>
        <Button style={{margin: '10px'}} type="primary" block onClick={()=>navigate("/kreis")}>Kreise</Button>
    </div>
  );
}

export default MainMenue;
