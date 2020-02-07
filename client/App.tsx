import React, { useEffect, useState, useCallback } from "react";
import ReactDOM from "react-dom";
import { BrowserRouter, Route, Redirect, Switch } from "react-router-dom";
import StatusPage from "./pages/StatusPage";
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
import { TransitionGroup, CSSTransition } from "react-transition-group";
import UserPage from "./pages/UserPage";
import SettingsPage from "./pages/SettingsPage";
import SubmissionPage from "./pages/SubmissionPage";
import ContestPage from "./pages/ContestPage";

function App() {
  let [user, setUser] = useState<User | null>(null);
  let [loading, setLoading] = useState(true);

  let fetchUserInformation = useCallback(function() {
    return request("/api/user")
      .then(function(user: User) {
        setUser(user);
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
        return fetchUserInformation();
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
        <div className="page-body">
          <Route
            render={function({ location }) {
              return (
                <TransitionGroup component={null}>
                  <CSSTransition
                    key={location.key}
                    timeout={200}
                    classNames="fade">
                    <div className="switch-wrapper">
                      <Switch location={location}>
                        <Route
                          path="/"
                          exact
                          component={ProblemsetPage}></Route>
                        <Route path="/status" component={StatusPage}></Route>
                        <Route path="/login" component={LoginPage}></Route>
                        <Route
                          path="/register"
                          component={RegisterPage}></Route>

                        <Route
                          exact
                          path="/problem/:problemId"
                          render={props => (
                            <ProblemPage tab={0} {...props}></ProblemPage>
                          )}
                        />
                        <Route
                          path="/problem/:problemId/submit"
                          render={props => {
                            if (!user) {
                              return <Redirect to="/login" />;
                            }
                            return (
                              <ProblemPage tab={1} {...props}></ProblemPage>
                            );
                          }}
                        />
                        <Route
                          path="/problem/:problemId/my"
                          render={props => {
                            if (!user) {
                              return <Redirect to="/login" />;
                            }
                            return (
                              <ProblemPage tab={2} {...props}></ProblemPage>
                            );
                          }}
                        />
                        <Route
                          path="/problem/:problemId/status"
                          render={props => (
                            <ProblemPage tab={3} {...props}></ProblemPage>
                          )}
                        />

                        <Route
                          path="/user/:username"
                          component={UserPage}></Route>
                        <Route
                          path="/settings/profile"
                          render={function(props) {
                            return <SettingsPage tab={0} {...props} />;
                          }}></Route>
                        <Route
                          path="/settings/password"
                          render={function(props) {
                            return <SettingsPage tab={1} {...props} />;
                          }}></Route>

                        <Route
                          path="/contest"
                          exact
                          component={ContestPage}></Route>

                        <Route
                          path="/submission/:submissionId"
                          component={SubmissionPage}></Route>

                        <Route
                          path="/logout"
                          render={() => {
                            setUser(null);
                            clearToken();
                            return <Redirect to="/"></Redirect>;
                          }}></Route>
                      </Switch>
                    </div>
                  </CSSTransition>
                </TransitionGroup>
              );
            }}
          />
        </div>
        <Footer></Footer>
      </BrowserRouter>
    </AppContext.Provider>
  );
}

ReactDOM.render(<App />, document.getElementById("root"));
