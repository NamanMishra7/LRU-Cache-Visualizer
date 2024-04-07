import React from 'react';
import './App.css';

import Main from './components/Main';

const App = () => {
  return (
    <div className="App" style={{ width: "100%" }}>
      <h1>LRU Cache Visualizer</h1>
      <Main />
    </div>
  );
};

export default App;
