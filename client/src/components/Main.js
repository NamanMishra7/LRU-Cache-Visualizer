import React from "react";

import ControllPanel from "./ControllPanel";
import Cache from "./Cache";

const Main = () => {
  return (
    <div style={{ display: "flex", flexDirection: 'row', width: "100%", justifyContent: 'space-between', margin: "0 1rem" }}>
      <div style={{ width: "35%" }}>
        <ControllPanel />
      </div>
      <div style={{ width: "60%" }}>
        <Cache />
      </div>
    </div>
  );
}

export default Main;