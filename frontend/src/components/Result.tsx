import React from "react";

type ResultProps = {
  keyword: string;
  name: string;
};

const Result = ({ keyword, name }: ResultProps) => {
  const result = JSON.parse(keyword);

  const winner = result.winner;
  const feedback = result.feedback;

  return (
    <>
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
      {name == winner ? (
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
          marginBottom: "20px",
          textAlign: "center",
        }}
      >
        勝者：{winner}
      </h2>
      <p
        style={{
          fontSize: "1.2rem",
          color: "#666",
          lineHeight: "1.6",
          padding: "20px",
          backgroundColor: "#f5f5f5",
          borderRadius: "10px",
          maxWidth: "800px",
          margin: "0 auto",
        }}
      >
        {feedback}
      </p>
    </>
  );
};

export default Result;
