import React, { useCallback, FormEvent } from "react";
import request from "../helpers/request";

interface ChangePasswordFormElement extends HTMLFormElement {
  oldPassword: HTMLInputElement;
  newPassword: HTMLInputElement;
  repPassword: HTMLInputElement;
}

function ChangePasswordForm() {
  let onSubmit = useCallback(function(e: FormEvent) {
    e.preventDefault();
    let form = e.target as ChangePasswordFormElement;

    let oldPassword = form.oldPassword.value;
    let newPassword = form.newPassword.value;
    let repPassword = form.repPassword.value;

    if (newPassword !== repPassword) {
      alert("Password does not match");
      return;
    }

    request("/api/user/password", {
      method: "PUT",
      body: JSON.stringify({ oldPassword, newPassword }),
      headers: { "Content-Type": "application/json" }
    })
      .then(function() {
        alert("Update successfully");
      })
      .catch(function({ error }) {
        alert(error);
      });
  }, []);

  return (
    <form className="auth-form" onSubmit={onSubmit}>
      <div className="auth-form__header">
        <h1>Change password</h1>
        <div className="auth-form__header-icon">
          <i className="fa fa-cog"></i>
        </div>
      </div>
      <div className="auth-form__body">
        <label>
          <span>Current password</span>
          <input type="password" name="oldPassword" />
        </label>
        <label>
          <span>New password</span>
          <input type="password" name="newPassword" />
        </label>
        <label>
          <span>Retype password</span>
          <input type="password" name="repPassword" />
        </label>
        <button type="submit">Update</button>
      </div>
    </form>
  );
}

export default ChangePasswordForm;
