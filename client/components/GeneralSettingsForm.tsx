import React, { useCallback, FormEvent, useContext } from "react";
import request from "../helpers/request";
import { useHistory } from "react-router-dom";
import AppContext from "../contexts/AppContext";
import User from "../models/User";

interface InputFileElement extends HTMLInputElement {
  files: FileList;
}

interface SettingsFormElement extends HTMLFontElement {
  fullname: HTMLInputElement;
  file: InputFileElement;
}

function GeneralSettingsForm({ user }: { user: User }) {
  let username = user.username;
  let history = useHistory();
  let { setUser } = useContext(AppContext);

  let onSubmit = useCallback(
    function(e: FormEvent) {
      e.preventDefault();
      let form = e.target as SettingsFormElement;
      let promiseList = [];

      // Update non-file fields
      let fullname: string = form.fullname.value;
      promiseList.push(
        request("/api/c/user", {
          method: "PATCH",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ fullname })
        }).then(function(user: User) {
          delete user.profilePicture;
          setUser((prev: User | null) =>
            prev
              ? {
                  ...prev,
                  ...user
                }
              : null
          );
        })
      );

      // Update profile picture
      if (form.file.files.length > 0) {
        let body = new FormData();
        body.append("file", form.file.files[0]);
        promiseList.push(
          request("/api/c/user/profile-picture", {
            method: "POST",
            body
          })
            .then(function({ path }: { path: string }) {
              setUser((user: User | null) =>
                user
                  ? {
                      ...user,
                      profilePicture: path
                    }
                  : null
              );
            })
            .catch(console.log)
        );
      }

      // When both requests finish, redirect user to their user page
      Promise.all(promiseList).then(function() {
        history.push("/user/" + username);
      });
    },
    [username, history, setUser]
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
          <input
            type="text"
            name="fullname"
            defaultValue={user.fullname}></input>
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
