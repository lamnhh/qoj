import React, { useEffect, useState, useCallback } from "react";
import ReactDOM from "react-dom";
import { BrowserRouter, Route, Redirect, Switch } from "react-router-dom";
import SubmissionPage from "./pages/SubmissionPage";
import ProblemsetPage from "./pages/ProblemsetPage";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";
import request from "./helpers/request";
import { setAccessToken, clearToken } from "./helpers/auth";
import Header from "./components/Header";
import Footer from "./components/Footer";
import ProblemPage from "./pages/ProblemPage";
import User from "./models/User";
import AppContext from "./contexts/AppContext";
import "./styles/index.scss";

function App() {
  let [user, setUser] = useState<User | null>(null);
  let [loading, setLoading] = useState(true);

  let fetchUserInformation = useCallback(function() {
    request("/api/user")
      .then(function({ username, fullname }) {
        setUser({ username, fullname });
      })
      .catch(function() {
        setUser(null);
      });
  }, []);

  // Get access token upon start
  useEffect(function() {
    request("/api/refresh")
      .then(function({ accessToken }) {
        setAccessToken(accessToken);
        fetchUserInformation();
      })
      .catch(function() {})
      .then(function() {
        setLoading(false);
      });
  }, []);

  if (loading) {
    return null;
  }

  return (
    <AppContext.Provider
      value={{
        user,
        setUser,
        fetchUserInformation
      }}>
      <BrowserRouter>
        <Header></Header>
        <div>
          <Switch>
            <Route path="/" exact component={ProblemsetPage}></Route>
            <Route path="/status" component={SubmissionPage}></Route>
            <Route path="/login" component={LoginPage}></Route>
            <Route path="/register" component={RegisterPage}></Route>
            <Route path="/problem/:problemId" component={ProblemPage}></Route>
            <Route
              path="/logout"
              render={() => {
                setUser(null);
                clearToken();
                return <Redirect to="/"></Redirect>;
              }}></Route>
          </Switch>
        </div>
        <Footer></Footer>
      </BrowserRouter>
    </AppContext.Provider>
  );
}

ReactDOM.render(<App />, document.getElementById("root"));
