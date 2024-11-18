import React, { useState } from "react";

export const ActionPage = ({
  keyword,
  socket,
  name,
  clientId,
  roomLevel,
}: {
  keyword: string;
  socket: WebSocket | null;
  name: string;
  clientId: string;
  roomLevel: number;
}) => {
  const [answer, setAnswer] = useState("");
  const [isAnswered, setIsAnswered] = useState(false);
  const [countTime, setCountTime] = useState(10);
  function submitAnswer() {
    setIsAnswered(true);
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(
        JSON.stringify({
          type: "answer",
          data: { clientId, name, answer, keyword, countTime, roomLevel },
        })
      );
    } else {
      console.error("WebSocket is not connected");
    }
  }

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        gap: "20px",
        padding: "30px",
        maxWidth: "600px",
        margin: "0 auto",
      }}
    >
      <h1
        style={{
          fontSize: "2rem",
          color: "#333",
          marginBottom: "10px",
        }}
      >
        お題に対して韻を踏んでください
      </h1>
      <h2
        style={{
          fontSize: "1.8rem",
          color: "#2c5282",
          padding: "15px 30px",
          backgroundColor: "#ebf8ff",
          borderRadius: "10px",
          marginBottom: "20px",
        }}
      >
        お題: {keyword}
      </h2>
      <input
        type="text"
        placeholder="韻を踏んだ言葉を入力してください"
        value={answer}
        onChange={(e) => setAnswer(e.target.value)}
        style={{
          width: "80%",
          padding: "12px 20px",
          fontSize: "1.1rem",
          borderRadius: "8px",
          border: "2px solid #cbd5e0",
          outline: "none",
          transition: "border-color 0.3s ease",
          marginBottom: "15px",
        }}
      />
      <button
        type="submit"
        onClick={submitAnswer}
        disabled={isAnswered}
        style={{
          backgroundColor: isAnswered ? "#cbd5e0" : "#4299e1",
          color: "white",
          padding: "12px 30px",
          fontSize: "1.1rem",
          border: "none",
          borderRadius: "8px",
          cursor: isAnswered ? "not-allowed" : "pointer",
          transition: "background-color 0.3s ease",
        }}
      >
        回答する
      </button>
      {isAnswered && (
        <p
          style={{
            fontSize: "1.2rem",
            color: "#718096",
            marginTop: "20px",
          }}
        >
          相手の回答を待っています・・・
        </p>
      )}
    </div>
  );
};
