"use client";
import { useEffect, useRef, useState } from 'react';
import Countdown from './countdown/countdown';

export default function Home() {
  const [message, setMessage] = useState('待機中...');
  const socketRef = useRef<WebSocket | null>(null);
  const [start, setStart] = useState(false);
  
  useEffect(() => {
    // WebSocket 接続を確立
    const socket = new WebSocket('ws://localhost:8080/ws');
    socketRef.current = socket;

    socket.onopen = () => {
      console.log('Connected to server');
    };

    socket.onmessage = (event) => {
      console.log('Received message:', event.data);
      if (event.data === "start") {
        setStart(true);
        setMessage('二人のユーザーが接続しました！');
      }
    };

    socket.onclose = () => {
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
            start ? <Countdown />: <div>hogehoge</div>
        }
    </div>
  );
}
