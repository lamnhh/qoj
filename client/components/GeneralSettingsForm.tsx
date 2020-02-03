import React, { useCallback, FormEvent } from "react";
import request from "../helpers/request";
import { useHistory } from "react-router-dom";

interface InputFileElement extends HTMLInputElement {
  files: FileList;
}

interface SettingsFormElement extends HTMLFontElement {
  fullname: HTMLInputElement;
  file: InputFileElement;
}

function GeneralSettingsForm({ username }: { username: string }) {
  let history = useHistory();

  let onSubmit = useCallback(
    function(e: FormEvent) {
      e.preventDefault();
      let form = e.target as SettingsFormElement;

      // TODO: Update non-file fields

      // Update profile picture
      if (form.file.files.length > 0) {
        let body = new FormData();
        body.append("file", form.file.files[0]);
        request("/api/user/profile-picture", {
          method: "POST",
          body
        })
          .then(function() {
            history.push("/user/" + username);
          })
          .catch(console.log);
      }
    },
    [username, history]
  );

  return (
    <form className="auth-form" onSubmit={onSubmit}>
      <div className="auth-form__header">
        <h1>General</h1>
        <div className="auth-form__header-icon">
          <i className="fa fa-cog"></i>
        </div>
      </div>
      <div className="auth-form__body">
        <label>
          <span>Full name</span>
          <input type="text" name="fullname"></input>
        </label>
        <label>
          <span>Profile picture</span>
          <input type="file" name="file"></input>
        </label>
        <button type="submit">Update</button>
      </div>
    </form>
  );
}

export default GeneralSettingsForm;
