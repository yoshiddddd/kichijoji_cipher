// components/Countdown.js
import { useState, useEffect } from 'react';

export default function Countdown() {
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
          <h1>Countdown Complete!</h1>
        </div>
      ) : (
        <h1>Countdown: {count}</h1>
      )}
    </div>
  );
}
