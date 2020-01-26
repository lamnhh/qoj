import React from "react";

interface WSContextInterface {
  socket: WebSocket;
}

let warning: WSContextInterface = {
  get socket(): WebSocket {
    throw new Error("WSContext.Provider required");
  }
};

let WSContext = React.createContext<WSContextInterface>(warning);

export default WSContext;
