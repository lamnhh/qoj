import React, { useCallback, FormEvent, useContext } from "react";
import { setAccessToken } from "../helpers/auth";
import { useHistory } from "react-router-dom";
import AppContext from "../contexts/AppContext";

interface LoginFormElement extends HTMLFontElement {
  username: HTMLInputElement;
  password: HTMLInputElement;
}

function LoginForm() {
  let history = useHistory();
  let { fetchUserInformation } = useContext(AppContext);

  let handleLogin = useCallback(
    function(event: FormEvent) {
      event.preventDefault();

      let form = event.target as LoginFormElement;
      let username = form.username.value;
      let password = form.password.value;

      fetch("/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
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
          history.push("/");
        })
        .catch(function({ error }) {
          alert(error);
        });
    },
    [history, fetchUserInformation]
  );

  return (
    <form className="auth-form" onSubmit={handleLogin}>
      <div className="auth-form__header">
        <h1>Sign in</h1>
        <div className="auth-form__header-icon">
          <i className="fa fa-key"></i>
        </div>
      </div>
      <div className="auth-form__body">
        <label>
          <span>Username</span>
          <input type="text" name="username" required />
        </label>
        <label>
          <span>Password</span>
          <input type="password" name="password" required />
        </label>
        <button type="submit">Sign in</button>
      </div>
    </form>
  );
}

export default LoginForm;
