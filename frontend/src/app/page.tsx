"use client";
import { useEffect, useRef, useState } from "react";
import { Countdown } from "../components/countdown";
import DifficultyLevelButton from "@/components/DifficultyLevelButton";

export default function Home() {
  const [message, setMessage] = useState("ゲームを開始してください");
  const socketRef = useRef<WebSocket | null>(null);
  const [start, setStart] = useState(false);
  const [keyword, setKeyword] = useState("");
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [clientId, setClientId] = useState("");
  const [isGameStarted, setIsGameStarted] = useState(false);
  const [name, setName] = useState("");
  const [thinking, setThinking] = useState(false);
  const [result, setResult] = useState(false);
  const [roomLevel, setRoomLevel] = useState(0);

  const changeRoomLevel = (level: number) => {
    setRoomLevel(level);
  };

  const handleStartGame = () => {
    // WebSocket 接続を確立
    if (!name) {
      alert("名前を入力してください");
      return;
    }
    const newSocket = new WebSocket("ws://localhost:8080/ws");
    setSocket(newSocket);
    socketRef.current = newSocket;

    setIsGameStarted(true);
    newSocket.onopen = () => {
      console.log("Connected to server");
      setMessage("ユーザーを探しています...");
      newSocket.send(
        JSON.stringify({ type: "start", data: { roomLevel, name } })
      );
    };

    newSocket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setKeyword(data.word);
      setClientId(data.clientId);
      console.log(data);

      console.log("Received message:", event.data);
      if (data.signal === "start") {
        setStart(true);
        setMessage("二人のユーザーが接続しました！ゲーム開始！");
      }
      if (data.signal === "end") {
        setThinking(true);
      }
      if (data.signal === "result") {
        setThinking(false);
        setResult(true);
        newSocket.close();
      }
    };

    newSocket.onclose = () => {
      console.log("Disconnected from server");
      setMessage("接続が切れました");
    };
  };

  const handleReset = () => {
    window.location.reload();
  };
  useEffect(() => {
    // クリーンアップ: コンポーネントのアンマウント時にWebSocket接続を閉じる
    return () => {
      if (socketRef.current) {
        socketRef.current.close();
      }
    };
  }, []);

  return (
    <div style={{ textAlign: "center", marginTop: "50px" }}>
      {thinking ? (
        <h1>AIがジャッジしています...</h1>
      ) : result ? (
        <div>
          <h1>結果発表</h1>
          <h2>{keyword}</h2>
          <button onClick={handleReset}>もう一度プレイ</button>
        </div>
      ) : start ? (
        <Countdown
          keyword={keyword}
          socket={socket}
          name={name}
          clientId={clientId}
          roomLevel={roomLevel}
        />
      ) : (
        <div>
          <p>{message}</p>
          {!isGameStarted && (
            <div>
              <input
                type="text"
                placeholder="名前を入力してください"
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
              <div>
                <DifficultyLevelButton
                  label="初級"
                  roomLevel={roomLevel}
                  changeRoomLevel={changeRoomLevel}
                />
                <DifficultyLevelButton
                  label="中級"
                  roomLevel={roomLevel}
                  changeRoomLevel={changeRoomLevel}
                />
                <DifficultyLevelButton
                  label="上級"
                  roomLevel={roomLevel}
                  changeRoomLevel={changeRoomLevel}
                />
              </div>
              <button onClick={handleStartGame}>ゲームを開始</button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
