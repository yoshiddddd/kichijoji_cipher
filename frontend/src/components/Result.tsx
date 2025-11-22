import React from "react";

export type ResultData = {
  clientId: string;
  signal: string;
  word: string;
  winner: string;
  user1Name: string;
  user2Name: string;
  user1_answer: string;
  user2_answer: string;
  user1_point: number;
  user2_point: number;
  feedback: string;
};

type ResultProps = {
  name: string;
  gameResult: ResultData;
};

const Result = ({ name, gameResult }: ResultProps) => {
  if (!gameResult) return null;

  const {
    winner,
    feedback,
    user1Name,
    user2Name,
    user1_answer,
    user2_answer,
    user1_point,
    user2_point,
  } = gameResult;

  return (
    <div style={{ width: "100%", maxWidth: "800px", margin: "0 auto" }}>
      <h1
        style={{
          fontSize: "2.5rem",
          color: "#333",
          marginBottom: "30px",
          textAlign: "center",
        }}
      >
        結果発表
      </h1>

      {name === winner ? (
        <p
          style={{
            fontSize: "2rem",
            color: "#4CAF50",
            fontWeight: "bold",
            marginBottom: "20px",
            textAlign: "center",
          }}
        >
          あなたの勝利です！
        </p>
      ) : (
        <p
          style={{
            fontSize: "2rem",
            color: "#f44336",
            fontWeight: "bold",
            marginBottom: "20px",
            textAlign: "center",
          }}
        >
          あなたの敗北です...
        </p>
      )}

      <h2
        style={{
          fontSize: "1.8rem",
          color: "#2196F3",
          marginBottom: "30px",
          textAlign: "center",
        }}
      >
        勝者：{winner}
      </h2>

      <div
        style={{
          display: "flex",
          flexDirection: "column",
          gap: "20px",
          marginBottom: "30px",
        }}
      >
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            gap: "20px",
            flexWrap: "wrap",
          }}
        >
          {/* User 1 Card */}
          <div
            style={{
              flex: 1,
              minWidth: "300px",
              padding: "20px",
              border: "2px solid #ddd",
              borderRadius: "12px",
              backgroundColor: user1Name === winner ? "#e8f5e9" : "#fff",
              boxShadow: "0 2px 4px rgba(0,0,0,0.05)",
              borderColor: user1Name === winner ? "#4CAF50" : "#ddd",
            }}
          >
            <h3
              style={{
                fontSize: "1.5rem",
                marginBottom: "10px",
                color: "#333",
                borderBottom: "1px solid #eee",
                paddingBottom: "10px",
              }}
            >
              {user1Name}
            </h3>
            <p
              style={{
                fontSize: "1.2rem",
                color: "#666",
                marginBottom: "15px",
              }}
            >
              得点:{" "}
              <span
                style={{
                  fontWeight: "bold",
                  color: "#2196F3",
                  fontSize: "1.4rem",
                }}
              >
                {user1_point}
              </span>
            </p>
            <div>
              <p
                style={{
                  fontSize: "0.9rem",
                  color: "#888",
                  marginBottom: "5px",
                  fontWeight: "bold",
                }}
              >
                回答（韻）:
              </p>
              <p
                style={{
                  fontSize: "1.1rem",
                  color: "#333",
                  whiteSpace: "pre-wrap",
                  backgroundColor: "rgba(255,255,255,0.7)",
                  padding: "10px",
                  borderRadius: "6px",
                  border: "1px solid #eee",
                }}
              >
                {user1_answer}
              </p>
            </div>
          </div>

          {/* User 2 Card */}
          <div
            style={{
              flex: 1,
              minWidth: "300px",
              padding: "20px",
              border: "2px solid #ddd",
              borderRadius: "12px",
              backgroundColor: user2Name === winner ? "#e8f5e9" : "#fff",
              boxShadow: "0 2px 4px rgba(0,0,0,0.05)",
              borderColor: user2Name === winner ? "#4CAF50" : "#ddd",
            }}
          >
            <h3
              style={{
                fontSize: "1.5rem",
                marginBottom: "10px",
                color: "#333",
                borderBottom: "1px solid #eee",
                paddingBottom: "10px",
              }}
            >
              {user2Name}
            </h3>
            <p
              style={{
                fontSize: "1.2rem",
                color: "#666",
                marginBottom: "15px",
              }}
            >
              得点:{" "}
              <span
                style={{
                  fontWeight: "bold",
                  color: "#2196F3",
                  fontSize: "1.4rem",
                }}
              >
                {user2_point}
              </span>
            </p>
            <div>
              <p
                style={{
                  fontSize: "0.9rem",
                  color: "#888",
                  marginBottom: "5px",
                  fontWeight: "bold",
                }}
              >
                回答（韻）:
              </p>
              <p
                style={{
                  fontSize: "1.1rem",
                  color: "#333",
                  whiteSpace: "pre-wrap",
                  backgroundColor: "rgba(255,255,255,0.7)",
                  padding: "10px",
                  borderRadius: "6px",
                  border: "1px solid #eee",
                }}
              >
                {user2_answer}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div
        style={{
          padding: "25px",
          backgroundColor: "#f5f5f5",
          borderRadius: "12px",
          marginTop: "30px",
        }}
      >
        <h3
          style={{
            fontSize: "1.3rem",
            marginBottom: "15px",
            color: "#333",
            display: "flex",
            alignItems: "center",
            gap: "10px",
          }}
        >
          <span>🤖</span> AIからのフィードバック
        </h3>
        <p
          style={{
            fontSize: "1.1rem",
            color: "#555",
            lineHeight: "1.8",
            whiteSpace: "pre-wrap",
          }}
        >
          {feedback}
        </p>
      </div>
    </div>
  );
};

export default Result;
