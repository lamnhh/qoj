import React, { useCallback, FormEvent, useContext } from "react";
import { setAccessToken } from "../helpers/auth";
import { useHistory } from "react-router-dom";
import AppContext from "../contexts/AppContext";

interface RegisterFormElement extends HTMLFontElement {
  username: HTMLInputElement;
  fullname: HTMLInputElement;
  password: HTMLInputElement;
}

function RegisterForm() {
  let history = useHistory();
  let { fetchUserInformation } = useContext(AppContext);

  let handleLogin = useCallback(
    function(event: FormEvent) {
      event.preventDefault();
      let form = event.target as RegisterFormElement;
      let username = form.username.value;
      let fullname = form.fullname.value;
      let password = form.password.value;

      fetch("/api/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({ username, fullname, password })
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
    [history]
  );

  return (
    <form className="auth-form" onSubmit={handleLogin}>
      <div className="auth-form__header">
        <h1>Register</h1>
        <div className="auth-form__header-icon">
          <i className="fa fa-sign-in"></i>
        </div>
      </div>
      <div className="auth-form__body">
        <label>
          <span>Username</span>
          <input type="text" name="username" required />
        </label>
        <label>
          <span>Full name</span>
          <input type="text" name="fullname" required />
        </label>
        <label>
          <span>Password</span>
          <input type="password" name="password" required />
        </label>
        <button type="submit">Register</button>
      </div>
    </form>
  );
}

export default RegisterForm;
