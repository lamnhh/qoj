import React, { useState, useCallback } from "react";

function App() {
  let [count, setCount] = useState(0);

  let increase = useCallback(function() {
    setCount((p) => p + 1);
  }, []);

  return (
    <>
      <h1>Count: {count}</h1>
      <button type="button" onClick={increase}>
        +1
      </button>
    </>
  );
}

export default App;
