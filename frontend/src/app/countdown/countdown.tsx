// components/Countdown.js
import { useState, useEffect } from "react";
import { ActionPage } from "../action/action_page";

type CountdownProps = {
  keyword: string;
  socket: WebSocket | null;
  name: string;
  clientId: string;
};

export const Countdown = ({
  keyword,
  socket,
  name,
  clientId,
}: CountdownProps) => {
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
    <div style={{ textAlign: "center", marginTop: "50px" }}>
      {completed ? (
        <div>
          <ActionPage
            keyword={keyword}
            socket={socket}
            name={name}
            clientId={clientId}
          />
        </div>
      ) : (
        <h1>Countdown: {count}</h1>
      )}
    </div>
  );
};
