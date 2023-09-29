import React, { useState } from 'react';
import { Input, Button } from 'antd';

function Chat(props) {
  const [message, setMessage] = useState('');

  const handleSendMessage = () => {
      props.sendMessage(JSON.stringify({ mode: "msg", message: message }));
      setMessage('');
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter') {
      handleSendMessage();
    }
  };

  return (
    <div className="chat-window">
      <Input maxLength={150} autoFocus={true} placeholder='Nachricht' value={message} onChange={(e)=>setMessage(e.target.value)} onKeyUp={handleKeyPress}></Input>
      <Button className='chat-button' onClick={()=>handleSendMessage()}>Senden</Button>
      {props.chatMessages.slice(-5).map((msg, index) => (<p className='chat-msg' key={index}><b>{msg.name+" (" + msg.size + "): "}</b>{msg.message}</p>))}
    </div>
  );

}

export default React.memo(Chat);
