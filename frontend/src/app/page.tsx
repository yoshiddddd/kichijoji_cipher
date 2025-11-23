"use client";
import { useEffect, useRef, useState } from "react";
import Image from "next/image";
import { Countdown } from "../components/countdown";
import DifficultyLevelButton from "@/components/DifficultyLevelButton";
import Result, { ResultData } from "@/components/Result";
import { TutorialModal } from "@/components/TutorialModal";

export default function Home() {
  const [message, setMessage] = useState("レベルを選択してバトル準備");
  const socketRef = useRef<WebSocket | null>(null);
  const [start, setStart] = useState(false);
  const [keyword, setKeyword] = useState("");
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [clientId, setClientId] = useState("");
  const [isGameStarted, setIsGameStarted] = useState(false);
  const [showTutorial, setShowTutorial] = useState(true); // 初期表示でチュートリアルを出す

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
  const [gameResult, setGameResult] = useState<ResultData | null>(null);

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

        // data.word がJSON文字列になっているためパースする
        try {
          const innerData = JSON.parse(data.word);
          console.log("Parsed innerData:", innerData);

          // Difyのレスポンス形式（動的キー）に合わせて整形
          const formattedResult: ResultData = {
            clientId: data.clientId,
            signal: data.signal,
            word: data.word,
            winner: innerData.winner || "",
            user1Name: innerData.user1name || "",
            user2Name: innerData.user2name || "",
            // 動的キー（名前_answer）から値を取得
            user1_answer:
              innerData[`${innerData.user1name}_answer`] ||
              innerData.user1_answer ||
              "",
            user2_answer:
              innerData[`${innerData.user2name}_answer`] ||
              innerData.user2_answer ||
              "",
            // 動的キー（名前_point）から値を取得
            user1_point:
              innerData[`${innerData.user1name}_point`] ||
              innerData.user1_point ||
              0,
            user2_point:
              innerData[`${innerData.user2name}_point`] ||
              innerData.user2_point ||
              0,
            feedback: innerData.feedback || "",
          };

          setGameResult(formattedResult);
        } catch (e) {
          console.error("Failed to parse result word:", e);
          // パース失敗時はそのままdataを入れる（あるいはエラー表示）
          setGameResult(data);
        }

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

  const handleTutorialClose = () => {
    setShowTutorial(false);
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
        <div
          style={{
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            width: "100%",
            padding: "40px 20px",
            flex: 1,
            justifyContent: "center",
          }}
        >
          <div
            style={{
              maxWidth: "100%",
              width: "min(300px, 80vw)",
              marginBottom: "50px",
            }}
          >
            <Image
              src="/kichijoji-cipher-loading.webp"
              alt="AI Judging"
              width={800}
              height={600}
              style={{
                width: "100%",
                height: "auto",
                borderRadius: "12px",
                boxShadow: "0 4px 15px rgba(0,0,0,0.1)",
                animation: "sway 2s ease-in-out infinite",
              }}
              priority
            />
          </div>
          <h1
            style={{
              fontSize: "clamp(1.5rem, 5vw, 2.5rem)",
              textAlign: "center",
              wordBreak: "keep-all",
              animation: "pulse 2s infinite",
              marginTop: "20px",
            }}
          >
            AIがジャッジしています...
          </h1>
          <style jsx>{`
            @keyframes pulse {
              0% {
                opacity: 1;
              }
              50% {
                opacity: 0.6;
              }
              100% {
                opacity: 1;
              }
            }
            @keyframes sway {
              0% {
                transform: translateX(-10px);
              }
              50% {
                transform: translateX(10px);
              }
              100% {
                transform: translateX(-10px);
              }
            }
          `}</style>
        </div>
      ) : result && gameResult ? (
        <div>
          <Result name={name} gameResult={gameResult} />
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
          {showTutorial && <TutorialModal onClose={handleTutorialClose} />}
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
