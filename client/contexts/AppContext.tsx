import React from "react";
import User from "../models/User";
import DispatchWrapper from "../models/DispatchWrapper";

interface AppContextInterface {
  user: User | null;
  setUser: DispatchWrapper<User | null>;
  fetchUserInformation: () => void;
}

let warning: AppContextInterface = {
  get user(): User {
    throw new Error("AppContext.Provider required");
  },
  get setUser(): React.Dispatch<React.SetStateAction<User | null>> {
    throw new Error("AppContext.Provider required");
  },
  get fetchUserInformation(): () => void {
    throw new Error("AppContext.Provider required");
  }
};

let AppContext = React.createContext<AppContextInterface>(warning);

export default AppContext;
