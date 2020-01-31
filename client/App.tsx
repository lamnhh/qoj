import React, { useEffect } from "react";
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

function App() {
  // Get access token upon start
  useEffect(function() {
    request("/api/refresh")
      .then(function({ accessToken }) {
        setAccessToken(accessToken);
      })
      .catch(() => {});
  }, []);

  return (
    <BrowserRouter>
      <Header></Header>
      <Switch>
        <Route path="/" exact component={ProblemsetPage}></Route>
        <Route path="/status" component={SubmissionPage}></Route>
        <Route path="/login" component={LoginPage}></Route>
        <Route path="/register" component={RegisterPage}></Route>
        <Route path="/problem/:problemId" component={ProblemPage}></Route>
        <Route
          path="/logout"
          render={() => {
            clearToken();
            return <Redirect to="/"></Redirect>;
          }}></Route>
      </Switch>
      <Footer></Footer>
    </BrowserRouter>
  );
}

ReactDOM.render(<App />, document.getElementById("root"));
