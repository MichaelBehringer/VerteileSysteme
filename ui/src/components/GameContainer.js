import { useParams } from 'react-router';
import { baseUrl } from '../helper/RequestHelper';
import MyCircleSocket from './MyCircleSocket';

function GameContainer() {
  const { id } = useParams();
  return (<MyCircleSocket serverUrl={"ws://"+baseUrl+"/game/"+id+"/ws"}/>);
}

export default GameContainer;
