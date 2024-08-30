import React, { useEffect, useState } from 'react';
import './Chat.css';
import Message from '../components/Message/Message';

const Chat: React.FC = () => {
  const [messages, setMessages] = useState<{ message: string, role: string }[]>([]);
  const [input, setInput] = useState('');

    useEffect(() => {
        setMessages([
            {
                role: "Assistant",
                message: "Hello, what can I do for you?"
            },
            {
                role: "User",
                message: "i need a text about lore ipsum"
            },
            {
                role: "Assistant",
                message: "here it is: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam ac sagittis nibh. Curabitur suscipit interdum turpis id fringilla. Pellentesque eu fringilla dui. Ut nunc sapien, lacinia et viverra vitae, feugiat vel metus. Morbi eget lorem eros. Vestibulum sed urna metus. Nam at ligula iaculis, laoreet lectus in, pellentesque neque. Vivamus vel erat non velit convallis eleifend id vitae sapien. Ut venenatis, nisi vitae malesuada porttitor, turpis dui finibus diam, vel finibus turpis ipsum ut libero. Nulla placerat quis libero at consectetur. Integer scelerisque nunc in orci fermentum tempus.'"
            }
        ])
    },[]);


  const sendMessage = () => {
    if (input.trim()) {
      setMessages([...messages, { message: input, role: "User" }]);
      setInput('');

      setTimeout(() => {
        setMessages(prevMessages => [...prevMessages, { message: 'AI Response', role: "AI" }]);
      }, 1000);
    }
  };

  return (
    <div className="chat-container">
      <div className="chat-messages">
        {messages.map((message, index) => (
          <Message 
            key={index} 
            message={message.message} 
            isUser={message.role == "User"} 
          />
        ))}
      </div>
      <div className="chat-input-container">
        <input 
          type="text" 
          value={input} 
          onChange={(e) => setInput(e.target.value)} 
          onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
          placeholder="Type your message here..." 
          className="chat-input"
        />
        <button onClick={sendMessage} className="chat-send-button">
          Send
        </button>
      </div>
    </div>
  );
};

export default Chat;
