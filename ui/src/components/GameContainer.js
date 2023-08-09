import { useEffect, useState } from 'react';
import { useParams } from 'react-router';
import { doGetRequestParam } from '../helper/RequestHelper';
import { WebSocketDemo } from './WebSocketDemo';
import MyCircleSocket from './MyCircleSocket';

function GameContainer() {
  const { id } = useParams();
  const [serverUrl, setServerUrl] = useState();
  useEffect(() => {
    doGetRequestParam('getUrl', id).then(
      res => {
        setServerUrl(res.data)
      }
    )
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);
  return (serverUrl?true?<MyCircleSocket serverUrl={serverUrl}/>:<WebSocketDemo serverUrl={serverUrl}/>:<h1>loading</h1>
  );
}

export default GameContainer;
