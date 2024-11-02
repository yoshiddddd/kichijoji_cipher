// components/Countdown.js
import { useState, useEffect } from 'react';
import { ActionPage } from '../action/action_page';

export const Countdown = ({keyword, socket, clientId}:{keyword:string,socket:WebSocket|null,clientId:string}) =>{
  const [count, setCount] = useState(3);
  const [completed, setCompleted] = useState(false);

  useEffect(() => {
    const intervalId = setInterval(() => {
      setCount((prevCount) => {
        if (prevCount <= 1) {
          clearInterval(intervalId);
          setCompleted(true); // カウントダウン完了後に切り替え
          return 0;
        }
        return prevCount - 1;
      });
    }, 1000);

    return () => clearInterval(intervalId); // クリーンアップ
  }, []);

  return (
    <div style={{ textAlign: 'center', marginTop: '50px' }}>
      {completed ? (
        <div>
          <ActionPage keyword={keyword} socket={socket} clientId={clientId}/>
        </div>
      ) : (
        <h1>Countdown: {count}</h1>
      )}
    </div>
  );
}
