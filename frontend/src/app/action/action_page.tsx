

export const ActionPage = ({keyword,socket}:{keyword:string,socket:WebSocket|null}) => {

    function submitAnswer() {
        console.log("submitAnswer");
        // socket?.send("submitAnswer");
    }

  return (

    <div>
      <h1>お題に対して韻を踏んでください</h1>
        <h2>お題: {keyword}</h2>
      <input type="text" placeholder="Enter your name" />
      <button onSubmit={submitAnswer} >Submit</button>
    </div>
  );
}