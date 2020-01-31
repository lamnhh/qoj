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
      }).then(function(res) {
        res.json().then(function({ accessToken }) {
          setAccessToken(accessToken);
          fetchUserInformation();
          history.push("/");
        });
      });
    },
    [history]
  );

  return (
    <form onSubmit={handleLogin}>
      <input type="text" name="username" placeholder="Username" required />
      <input type="text" name="fullname" placeholder="Full name" required />
      <input type="password" name="password" placeholder="Password" required />
      <button type="submit">Register</button>
    </form>
  );
}

export default RegisterForm;
