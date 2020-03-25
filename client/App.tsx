import React, { useEffect, useState, useCallback } from "react";
import ReactDOM from "react-dom";
import { BrowserRouter, Route, Redirect, Switch } from "react-router-dom";
import { TransitionGroup, CSSTransition } from "react-transition-group";
import loadable from "@loadable/component";

import request from "./helpers/request";
import { setAccessToken } from "./helpers/auth";
import Header from "./components/Header";
import Footer from "./components/Footer";
import User from "./models/User";
import AppContext from "./contexts/AppContext";
import "./styles/index.scss";
import Logout from "./components/Logout";

const StatusPage = loadable(() => import("./pages/StatusPage"), {});
const ProblemsetPage = loadable(() => import("./pages/ProblemsetPage"));
const LoginPage = loadable(() => import("./pages/LoginPage"));
const RegisterPage = loadable(() => import("./pages/RegisterPage"));
const ProblemPage = loadable(() => import("./pages/ProblemPage"));
const UserPage = loadable(() => import("./pages/UserPage"));
const SettingsPage = loadable(() => import("./pages/SettingsPage"));
const SubmissionPage = loadable(() => import("./pages/SubmissionPage"));
const ContestListPage = loadable(() => import("./pages/ContestListPage"));
const ContestParticipantsPage = loadable(() =>
  import("./pages/ContestParticipantsPage")
);
const ContestPage = loadable(() => import("./pages/ContestPage"));

// Initialise moment-duration
let moment = require("moment");
let momentDuration = require("moment-duration-format");
momentDuration(moment);

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
  useEffect(
    function() {
      request("/api/refresh")
        .then(function({ accessToken }) {
          setAccessToken(accessToken);
          return fetchUserInformation();
        })
        .finally(function() {
          setLoading(false);
        });
    },
    [fetchUserInformation]
  );

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
                          component={ContestListPage}></Route>
                        <Route
                          path="/contest/:contestId/participants"
                          component={ContestParticipantsPage}></Route>

                        <Route
                          path="/contest/:contestId"
                          exact
                          render={props => (
                            <ContestPage {...props} tab={0} />
                          )}></Route>
                        <Route
                          path="/contest/:contestId/submit"
                          render={props => {
                            if (!user) {
                              return <Redirect to="/login" />;
                            }
                            return <ContestPage {...props} tab={1} />;
                          }}></Route>
                        <Route
                          path="/contest/:contestId/my"
                          render={props => {
                            if (!user) {
                              return <Redirect to="/login" />;
                            }
                            return <ContestPage {...props} tab={2} />;
                          }}></Route>
                        <Route
                          path="/contest/:contestId/status"
                          render={props => (
                            <ContestPage {...props} tab={3} />
                          )}></Route>

                        <Route
                          path="/contest/:contestId/ranking"
                          render={props => (
                            <ContestPage {...props} tab={4} />
                          )}></Route>

                        <Route
                          path="/submission/:submissionId"
                          component={SubmissionPage}></Route>

                        <Route path="/logout">
                          <Logout logoutUrl="/api/logout" returnUrl="/" />
                        </Route>
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
