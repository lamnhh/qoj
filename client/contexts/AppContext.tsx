import React, { Dispatch } from "react";
import User from "../models/User";

type UserNull = User | null;

type UserAction = UserNull | ((prevState: UserNull) => UserNull);

interface AppContextInterface {
  user: User | null;
  setUser: Dispatch<UserAction>;
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
