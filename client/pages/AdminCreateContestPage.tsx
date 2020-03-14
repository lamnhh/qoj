import React from "react";
import AdminContestForm from "../components/AdminContestForm";
import { defaultContest } from "../models/Contest";
import AdminContest from "../models/AdminContest";
import request from "../helpers/request";
import { useHistory } from "react-router-dom";

function AdminCreateContestPage() {
  let history = useHistory();

  function onCreate(contest: AdminContest) {
    request("/api/contest", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(contest)
    })
      .then(function() {
        history.push("/contest");
      })
      .catch(function({ error }) {
        alert(error);
      });
  }

  return (
    <AdminContestForm
      action="Create"
      defaultContest={defaultContest}
      defaultSelected={[]}
      handleSubmit={onCreate}
    />
  );
}

export default AdminCreateContestPage;
