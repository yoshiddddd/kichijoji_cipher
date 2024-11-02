"use client";
import { useEffect, useRef, useState } from 'react';
import {Countdown} from './countdown/countdown';

export default function Home() {
  const [message, setMessage] = useState('待機中...');
  const socketRef = useRef<WebSocket | null>(null);
  const [start, setStart] = useState(false);
  const [keyword, setKeyword] = useState('');
  const [socket, setSocket] = useState<WebSocket |null>(null);
  const [clientId, setClientId] = useState('');
  useEffect(() => {
    // WebSocket 接続を確立
    const newSocket = new WebSocket('ws://localhost:8080/ws');
    setSocket(newSocket);
    socketRef.current = socket;
        newSocket.onopen = () => {
          console.log('Connected to server');
        };
    
        newSocket.onmessage = (event) => {
         const data = JSON.parse(event.data);
            setKeyword(data.word);
            setClientId(data.clientId);
            console.log(data);
    
          console.log('Received message:', event.data);
          if (data.signal === "start") {
            setStart(true);
            setMessage('二人のユーザーが接続しました！');
          }
        };
    
        newSocket.onclose = () => {
          console.log('Disconnected from server');
        };
    
        return () => {
          if (socketRef.current) {
            socketRef.current.close();
          }
        };
  }, []);

  return (
    <div style={{ textAlign: 'center', marginTop: '50px' }}>
        {
            start ? <Countdown keyword ={keyword} socket={socket} clientId={clientId} />: <div>hogehoge</div>
        }
    </div>
  );
}
