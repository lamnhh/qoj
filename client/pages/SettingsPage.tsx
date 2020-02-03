import React, { useContext } from "react";
import GeneralSettingsForm from "../components/GeneralSettingsForm";
import AppContext from "../contexts/AppContext";
import { Redirect } from "react-router-dom";

function SettingsPage() {
  let user = useContext(AppContext).user;
  if (user === null) {
    return <Redirect to="/login"></Redirect>;
  }

  return (
    <>
      <header className="page-name align-left-right">
        <h1>Settings</h1>
      </header>
      <section className="settings-page align-left-right">
        <GeneralSettingsForm username={user.username} />
      </section>
    </>
  );
}

export default SettingsPage;
