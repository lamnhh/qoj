import React, { useContext } from "react";
import GeneralSettingsForm from "../components/GeneralSettingsForm";
import AppContext from "../contexts/AppContext";
import { Redirect } from "react-router-dom";
import { Tab, Tabs, TabList, TabPanel } from "react-tabs";

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
      <Tabs className="settings-page align-left-right">
        <TabList className="tab-list">
          <Tab>General Settings</Tab>
          <Tab>Change Password</Tab>
        </TabList>
        <TabPanel>
          <GeneralSettingsForm user={user} />
          <div>Change Password</div>
        </TabPanel>
      </Tabs>
    </>
  );
}

export default SettingsPage;
