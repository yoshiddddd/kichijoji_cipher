

export const ActionPage = ({keyword}:{keyword:string}) => {
    
  return (

    <div>
      <h1>お題に対して韻を踏んでください</h1>
        <h2>お題: {keyword}</h2>
      <input type="text" placeholder="Enter your name" />
      <button  >Submit</button>
    </div>
  );
}