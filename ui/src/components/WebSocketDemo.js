import { Button, Input, Table } from 'antd';
import React, { useState, useEffect } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';

export const WebSocketDemo = (props) => {
  const [myMessageHistory, setMyMessageHistory] = useState([]);
  const [inputText, setInputText] = useState("");

  const { sendMessage, lastMessage, readyState } = useWebSocket(props.serverUrl);

  useEffect(() => {
    if (lastMessage !== null) {
      const newMessage = {
        key: lastMessage.timeStamp,
        sender: lastMessage.data.substring(0,36),
        text: lastMessage.data.substring(38)
      }
      setMyMessageHistory(myMessageHistory => [...myMessageHistory, newMessage])
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [lastMessage]);

  function myHandleClickSendMessage() {
    sendMessage(inputText)
    setInputText("")
  }

  const connectionStatus = {
    [ReadyState.CONNECTING]: 'Connecting ↩️',
    [ReadyState.OPEN]: 'Verbunden 👍',
    [ReadyState.CLOSING]: 'Geschlossen 👎',
    [ReadyState.CLOSED]: 'Geschlossen 👎',
    [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
  }[readyState];

  const columns = [
    {
      title: 'Sender',
      dataIndex: 'sender',
      key: 'sender',
    },
    {
      title: 'Text',
      dataIndex: 'text',
      key: 'text',
    }
  ];

  return (
    <div>
      <h1>Status: {connectionStatus}</h1>
      <Input value={inputText} onChange={(e)=>setInputText(e.target.value)} placeholder='Eingabetext'/>
      <Button onClick={()=>myHandleClickSendMessage()} block type='primary'>Senden</Button>
      <Table columns={columns} dataSource={myMessageHistory} />
    </div>
  );
};