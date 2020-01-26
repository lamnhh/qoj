import React, { useCallback, FormEvent } from "react";
import { setAccessToken } from "../helpers/auth";

interface RegisterFormElement extends HTMLFontElement {
  username: HTMLInputElement;
  fullname: HTMLInputElement;
  password: HTMLInputElement;
}

function RegisterForm() {
  let handleLogin = useCallback(function(event: FormEvent) {
    event.preventDefault();
    let form = event.target as RegisterFormElement;
    let username = form.username.value;
    let fullname = form.fullname.value;
    let password = form.password.value;

    fetch("/api/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ username, fullname, password })
    }).then(function(res) {
      res.json().then(function({ accessToken }) {
        setAccessToken(accessToken);
      });
    });
  }, []);

  return (
    <form onSubmit={handleLogin}>
      <input type="text" name="username"></input>
      <input type="text" name="fullname"></input>
      <input type="password" name="password"></input>
      <button type="submit">Login</button>
    </form>
  );
}

export default RegisterForm;
