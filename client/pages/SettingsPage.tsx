import React, { useContext, useCallback } from "react";
import GeneralSettingsForm from "../components/GeneralSettingsForm";
import AppContext from "../contexts/AppContext";
import { Redirect, useHistory } from "react-router-dom";
import { Tab, Tabs, TabList, TabPanel } from "react-tabs";
import ChangePasswordForm from "../components/ChangePasswordForm";

function SettingsPage({ tab }: { tab: number }) {
  let user = useContext(AppContext).user;
  let history = useHistory();

  let onTabChange = useCallback(
    function(tab) {
      if (tab === 0) {
        history.push("/settings/profile");
      } else {
        history.push("/settings/password");
      }
    },
    [history]
  );

  if (user === null) {
    return <Redirect to="/login"></Redirect>;
  }
  return (
    <>
      <header className="page-name align-left-right">
        <h1>Settings</h1>
      </header>
      <Tabs
        className="settings-page align-left-right"
        selectedIndex={tab}
        onSelect={onTabChange}>
        <TabList className="my-tablist">
          <Tab className="my-tab">General Settings</Tab>
          <Tab className="my-tab">Change Password</Tab>
        </TabList>
        <TabPanel>
          <GeneralSettingsForm user={user} />
        </TabPanel>
        <TabPanel>
          <ChangePasswordForm />
        </TabPanel>
      </Tabs>
    </>
  );
}

export default SettingsPage;
