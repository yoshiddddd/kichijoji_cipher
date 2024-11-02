import React, { useState } from "react";

export const ActionPage = ({keyword,socket,clientId}:{keyword:string,socket:WebSocket|null,clientId:string}) => {
    const [answer, setAnswer] = useState("");
    const [isAnswered, setIsAnswered] = useState(false);
    function submitAnswer() {
        setIsAnswered(true);
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify({ clientId, answer, keyword }));
            setAnswer(''); // 入力フィールドをリセット
          } else {
            console.error('WebSocket is not connected');
          }
    }

    return (
        <div>
          <h1>お題に対して韻を踏んでください</h1>
          <h2>お題: {keyword}</h2>
            <input
              type="text"
              placeholder="Enter your answer"
              value={answer}
                onChange={(e) => setAnswer(e.target.value)}
            //   onChange={(e) => setAnswer(e.target.value)}
            />
            <button type="submit" onClick={submitAnswer} disabled={isAnswered}>Submit</button>
        </div>
      );
}