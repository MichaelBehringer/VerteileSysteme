import { useParams } from 'react-router';
import { baseUrl } from '../helper/RequestHelper';
import WebsocketContainer from './WebsocketContainer';

function GameContainer(props) {
  const { id } = useParams();
  if(id==="random") {
    return (<WebsocketContainer token={props.token} serverUrl={"ws://"+baseUrl+"/randGame/ws"}/>);
  } else {
    return (<WebsocketContainer token={props.token} serverUrl={"ws://"+baseUrl+"/game/"+id+"/ws"}/>);
  }
}

export default GameContainer;
