import React, { useCallback, FormEvent } from "react";
import { setAccessToken } from "../helpers/auth";

interface LoginFormElement extends HTMLFontElement {
  username: HTMLInputElement;
  password: HTMLInputElement;
}

function LoginForm() {
  let handleLogin = useCallback(function(event: FormEvent) {
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
    }).then(function(res) {
      res.json().then(function({ accessToken }) {
        setAccessToken(accessToken);
      });
    });
  }, []);

  return (
    <form onSubmit={handleLogin}>
      <input type="text" name="username"></input>
      <input type="password" name="password"></input>
      <button type="submit">Login</button>
    </form>
  );
}

export default LoginForm;
