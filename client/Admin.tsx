import React, { useState, useCallback, useEffect } from "react";
import ReactDOM from "react-dom";
import "./styles/admin.scss";
import { BrowserRouter, Route, Switch } from "react-router-dom";
import AdminLayout from "./components/AdminLayout";
import AdminLoginPage from "./pages/AdminLoginPage";
import User from "./models/User";
import request from "./helpers/request";
import { setAccessToken } from "./helpers/auth";
import AppContext from "./contexts/AppContext";
import Logout from "./components/Logout";
import AdminCreateProblemPage from "./pages/AdminCreateProblemPage";
import AdminProblemPage from "./pages/AdminProblemPage";
import AdminEditProblemPage from "./pages/AdminEditProblemPage";
import AdminSearchPage from "./pages/AdminSearchPage";
import AdminContestPage from "./pages/AdminContestPage";

// Initialise moment-duration
let moment = require("moment");
let momentDuration = require("moment-duration-format");
momentDuration(moment);

function Admin() {
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
        <Switch>
          <Route path="/" exact>
            <AdminLayout title="Login" requireLogin={false}>
              <AdminLoginPage />
            </AdminLayout>
          </Route>
          <Route path="/logout">
            <Logout logoutUrl="/api/logout" returnUrl="/" />
          </Route>
          <Route path="/search">
            <AdminLayout title="Search">
              <AdminSearchPage />
            </AdminLayout>
          </Route>
          <Route path="/problem/new">
            <AdminLayout title="Create Problem">
              <AdminCreateProblemPage />
            </AdminLayout>
          </Route>
          <Route path="/problem" exact>
            <AdminLayout title="Problems">
              <AdminProblemPage />
            </AdminLayout>
          </Route>
          <Route path="/problem/edit/:id">
            <AdminLayout title="Edit problem">
              <AdminEditProblemPage />
            </AdminLayout>
          </Route>
          <Route path="/contest/new">
            <AdminLayout title="Create Contest">Create contest</AdminLayout>
          </Route>
          <Route path="/contest" exact>
            <AdminLayout title="Contests">
              <AdminContestPage />
            </AdminLayout>
          </Route>
        </Switch>
      </BrowserRouter>
    </AppContext.Provider>
  );
}

ReactDOM.render(<Admin />, document.getElementById("root"));
