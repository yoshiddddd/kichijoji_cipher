"use client";
import { useEffect, useRef, useState } from "react";
import { Countdown } from "../components/countdown";
import DifficultyLevelButton from "@/components/DifficultyLevelButton";
import Result from "@/components/Result";

export default function Home() {
  const [message, setMessage] = useState("レベルを選択してバトル準備");
  const socketRef = useRef<WebSocket | null>(null);
  const [start, setStart] = useState(false);
  const [keyword, setKeyword] = useState("");
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [clientId, setClientId] = useState("");
  const [isGameStarted, setIsGameStarted] = useState(false);
  const [name, setName] = useState(() => {
    if (typeof window !== "undefined") {
      const savedName = localStorage.getItem("userName");
      return savedName || "";
    }
    return "";
  });
  const [thinking, setThinking] = useState(false);
  const [result, setResult] = useState(false);
  const [roomLevel, setRoomLevel] = useState(0);
  const [secretWord, setSecretWord] = useState("");

  const changeRoomLevel = (level: number) => {
    setRoomLevel(level);
  };

  const handleStartGame = () => {
    const newSocket = new WebSocket(
      "wss://kichijoji-cipher-be-deploy.onrender.com/ws"
    );
    setSocket(newSocket);
    socketRef.current = newSocket;

    setIsGameStarted(true);
    newSocket.onopen = () => {
      console.log("Connected to server");
      setMessage("ユーザーを探しています...");
      newSocket.send(
        JSON.stringify({ type: "start", data: { roomLevel, name, secretWord } })
      );
    };

    newSocket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setKeyword(data.word);
      setClientId(data.clientId);
      console.log(data);

      console.log("Received message:", event.data);
      if (data.signal === "userLeft") {
        alert("相手が退出しました.");
        window.location.reload();
      }
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
      if (data.signal === "alreadyExist") {
        alert("すでに同じ合言葉が存在します。名前を変更してください。");
        window.location.reload();
      }
    };

    newSocket.onclose = () => {
      console.log("Disconnected from server");
      setMessage("接続が切れました");
    };

    // 名前をローカルストレージに保存
    if (typeof window !== "undefined" && name) {
      localStorage.setItem("userName", name);
    }

    // WebSocket 接続を確立
    if (!name) {
      alert("名前を入力してください");
      return;
    }
    if (!secretWord) {
      alert("合言葉を入力してください");
      return;
    }
    if (!roomLevel) {
      alert("レベルを選択してください");
      return;
    }
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
    <div
      style={{
        textAlign: "center",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        height: "100vh",
      }}
    >
      {thinking ? (
        <h1>AIがジャッジしています...</h1>
      ) : result ? (
        <div>
          <Result keyword={keyword} name={name} />
          <button
            onClick={handleReset}
            style={{
              backgroundColor: "#4299e1",
              color: "white",
              padding: "12px 30px",
              fontSize: "1.1rem",
              border: "none",
              borderRadius: "8px",
              cursor: "pointer",
              transition: "background-color 0.3s ease",
              marginTop: "30px",
            }}
          >
            もう一度プレイ
          </button>
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
          <p
            style={{ fontSize: "1.2rem", color: "black", marginBottom: "20px" }}
          >
            {message}
          </p>
          {!isGameStarted && (
            <div
              style={{
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                gap: "20px",
              }}
            >
              <input
                type="text"
                placeholder="名前を入力してください"
                value={name}
                onChange={(e) => setName(e.target.value)}
                style={{
                  padding: "12px 20px",
                  fontSize: "1rem",
                  borderRadius: "8px",
                  border: "2px solid #ddd",
                  width: "300px",
                  outline: "none",
                  transition: "border-color 0.3s ease",
                }}
              />
              <input
                type="text"
                placeholder="合言葉を入力してください。"
                value={secretWord}
                onChange={(e) => setSecretWord(e.target.value)}
                style={{
                  padding: "12px 20px",
                  fontSize: "1rem",
                  borderRadius: "8px",
                  border: "2px solid #ddd",
                  width: "300px",
                  outline: "none",
                  transition: "border-color 0.3s ease",
                }}
              />
              <div style={{ margin: "20px 0" }}>
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
              <button
                onClick={handleStartGame}
                style={{
                  backgroundColor: "#4CAF50",
                  color: "white",
                  padding: "15px 30px",
                  fontSize: "1.1rem",
                  border: "none",
                  borderRadius: "8px",
                  cursor: "pointer",
                  transition: "background-color 0.3s ease",
                }}
              >
                ゲームを開始
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
