import React, { useCallback, useEffect, useRef, useState } from 'react';
import './Chat.css';
import Message from '../components/Message/Message';
import { useAuth } from '@clerk/clerk-react';
import { useNavigate } from 'react-router';

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

const Chat: React.FC = () => {
  const [messages, setMessages] = useState<{ content: string, sender: string }[]>([]);
  const [eventId, setEventId] = useState<number | undefined>();
  const [input, setInput] = useState('');
  const { getToken } = useAuth();
  const [isLoading, setIsLoading] = useState(false);

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate();

  const sendMessage = async () => {
    try {
      const token = await getToken();

      setInput('');

      setMessages((prevMessages = []) => [
        ...prevMessages,
        { content: input, sender: "User" },
        { content: "AI is typing...", sender: "AI" }
      ]);
      setIsLoading(true);

      const response = await fetch(`${BACKEND_URL}/answer/${eventId}`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ answer: input }),
      });

      if (!response.ok) {
        throw new Error('Error sending message to server');
      }

      const responseData = await response.json();

      setMessages((prevMessages = []) => {
        const updatedMessages = [...prevMessages];
        updatedMessages.pop();  
        return [...updatedMessages, responseData.message];
      });

      if (responseData.message.content === "Ha terminado la entrevista. Puedes proceder a ver los resultados") {
        navigate('/results');
      }

    } catch (error) {
      alert(error);
    } finally {
      setIsLoading(false);
    }
  };

  const getFirstQuestion = useCallback(async () => {
    try {
      const token = await getToken();

      if(eventId == null){
        return;
      }

      const response = await fetch(`${BACKEND_URL}/questions/${eventId}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
      });

      if (!response.ok) {
        throw new Error('Error getting the question');
      }

      const responseData = await response.json();

      setMessages((prevMessages = []) => [...prevMessages, responseData.message]);

    } catch (error) {
      alert(error);
    }
  }, [getToken, eventId]);

  const startEvent = useCallback(async () => {
    try {
      const token = await getToken();

      const response = await fetch(`${BACKEND_URL}/start-test`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
      });

      if (!response.ok) {
        throw new Error('Error starting the career survey');
      }

      const responseData = await response.json();

      setEventId(responseData.careerTest.id);

      if (responseData.careerTest.fullConversation !== null) {
        return;
      }

      console.log("getting question");
      await getFirstQuestion();

    } catch (error) {
      alert(error);
    }
  }, [getToken, getFirstQuestion]);

  const fetchCareerTest = useCallback(async () => {
    try {
      const token = await getToken();

      const response = await fetch(`${BACKEND_URL}/get-test`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      if (response.status === 404) {
        startEvent();
        return;
      }

      if (!response.ok) {
        throw new Error('Error loading the career survey');
      }

      const responseData = await response.json();

      setMessages((responseData.careerTest.fullConversation) || []);
      setEventId(responseData.careerTest.id);

      if (responseData.careerTest.fullConversation !== null) {
        return;
      }

      await getFirstQuestion();
    } catch (error) {
      alert(error);
    }
  }, [getToken, startEvent, getFirstQuestion]);

  useEffect(() => {
    fetchCareerTest();
  }, [fetchCareerTest]);

  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages]);

  return (
    <div className="chat-container">
      <div className="chat-messages">
        {messages?.map((message, index) => (
          <Message
            key={index}
            message={message.content}
            isUser={message.sender === "User"}
          />
        ))}
        <div ref={messagesEndRef} />
      </div>
      <div className="chat-input-container">
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && sendMessage()}
          placeholder="Type your message here..."
          className="chat-input"
          disabled={isLoading}
        />
        <button onClick={sendMessage} className="chat-send-button" disabled={isLoading}>
          Send
        </button>
      </div>
    </div>
  );
};

export default Chat;
