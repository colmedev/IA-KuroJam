import React from 'react';
import './Message.css';

interface MessageProps {
  message: string;
  isUser: boolean;
}

const Message: React.FC<MessageProps> = ({ message, isUser }) => {
  return (
    <div className={`message ${isUser ? 'user-message' : 'ai-message'}`}>
      <div className="message-content">
        {message}
      </div>
    </div>
  );
};

export default Message;
