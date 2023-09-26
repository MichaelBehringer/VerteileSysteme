import { useParams } from 'react-router';
import { baseUrl } from '../helper/RequestHelper';
import MyCircleSocket from './MyCircleSocket';

function GameContainer(props) {
  const { id } = useParams();
  if(id==="random") {
    return (<MyCircleSocket token={props.token} serverUrl={"ws://"+baseUrl+"/randGame/ws"}/>);
  } else {
    return (<MyCircleSocket token={props.token} serverUrl={"ws://"+baseUrl+"/game/"+id+"/ws"}/>);
  }
}

export default GameContainer;
