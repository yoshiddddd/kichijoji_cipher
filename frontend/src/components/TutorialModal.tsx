import React from "react";

type TutorialModalProps = {
  onClose: () => void;
};

export const TutorialModal = ({ onClose }: TutorialModalProps) => {
  return (
    <div
      style={{
        position: "fixed",
        top: 0,
        left: 0,
        width: "100%",
        height: "100%",
        backgroundColor: "rgba(0, 0, 0, 0.5)",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        zIndex: 1000,
      }}
    >
      <div
        style={{
          backgroundColor: "white",
          padding: "30px",
          borderRadius: "12px",
          maxWidth: "500px",
          width: "90%",
          boxShadow: "0 4px 15px rgba(0, 0, 0, 0.2)",
          textAlign: "center",
        }}
      >
        <h2
          style={{
            fontSize: "1.8rem",
            color: "#333",
            marginBottom: "20px",
            borderBottom: "2px solid #4CAF50",
            display: "inline-block",
            paddingBottom: "5px",
          }}
        >
          遊び方
        </h2>
        <div
          style={{
            textAlign: "left",
            marginBottom: "30px",
            lineHeight: "1.6",
            color: "#555",
            fontSize: "1.1rem",
          }}
        >
          <p style={{ marginBottom: "15px" }}>
            お題に対して<span style={{ fontWeight: "bold", color: "#E91E63" }}>韻を踏んだ単語（文）</span>を回答しましょう！
          </p>
          <p style={{ marginBottom: "15px" }}>
            <span style={{ fontWeight: "bold", color: "#2196F3" }}>母音の一致度</span>や
            <span style={{ fontWeight: "bold", color: "#2196F3" }}>文脈</span>、
            <span style={{ fontWeight: "bold", color: "#2196F3" }}>回答時間</span>をもとにAIが得点を計算して勝敗をジャッジします。
          </p>
        </div>
        <button
          onClick={onClose}
          style={{
            backgroundColor: "#4CAF50",
            color: "white",
            padding: "12px 40px",
            fontSize: "1.2rem",
            border: "none",
            borderRadius: "25px",
            cursor: "pointer",
            transition: "background-color 0.3s ease",
            fontWeight: "bold",
            boxShadow: "0 2px 5px rgba(0,0,0,0.2)",
          }}
          onMouseOver={(e) => (e.currentTarget.style.backgroundColor = "#45a049")}
          onMouseOut={(e) => (e.currentTarget.style.backgroundColor = "#4CAF50")}
        >
          OK
        </button>
      </div>
    </div>
  );
};
