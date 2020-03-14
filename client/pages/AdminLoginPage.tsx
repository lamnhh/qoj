import React, { FormEvent, useContext } from "react";
import { setAccessToken } from "../helpers/auth";
import { useHistory, Redirect } from "react-router-dom";
import AppContext from "../contexts/AppContext";

interface LoginFormElement extends HTMLFormElement {
  username: HTMLInputElement;
  password: HTMLInputElement;
}

function AdminLoginPage() {
  let { user, fetchUserInformation } = useContext(AppContext);
  let history = useHistory();

  function onLogin(e: FormEvent) {
    e.preventDefault();
    let form = e.target as LoginFormElement;

    let username = form.username.value;
    let password = form.password.value;

    fetch("/api/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password })
    })
      .then(function(res) {
        return res.json().then(function(data) {
          if (res.ok) {
            return data;
          }
          throw data;
        });
      })
      .then(function({ accessToken }) {
        setAccessToken(accessToken);
        fetchUserInformation();
        history.push("/problem");
      })
      .catch(function({ error }) {
        alert(error);
      });
  }

  if (user) {
    return <Redirect to="/problem" />;
  }

  return (
    <section className="login-page--wrapper">
      <div className="login-page">
        <h2>Login Form</h2>
        <form onSubmit={onLogin}>
          <label>
            <span>Username:</span>
            <input type="text" name="username" required autoFocus />
          </label>
          <label>
            <span>Password:</span>
            <input type="password" name="password" required />
          </label>
          <button type="submit">Login</button>
        </form>
      </div>
    </section>
  );
}

export default AdminLoginPage;
