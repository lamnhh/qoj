import React from "react";
import AdminContestForm from "../components/AdminContestForm";
import { useParams, useHistory } from "react-router-dom";
import useSWR from "swr";
import request from "../helpers/request";
import AdminContest from "../models/AdminContest";
import Contest from "../models/Contest";
import Loading from "../components/Loading";

interface AdminEditContestPageRouterProps {
  id: string;
}

function AdminEditContestPage() {
  let contestId = parseInt(useParams<AdminEditContestPageRouterProps>().id);
  let { data } = useSWR("/api/contest/" + contestId, request);

  let history = useHistory();

  function onUpdate(contest: AdminContest) {
    request("/api/contest/" + contestId, {
      method: "PATCH",
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

  if (!data) {
    return <Loading />;
  }

  let contest = data as Contest;
  let selected = (data as AdminContest).problemList;

  return (
    <AdminContestForm
      action="Update"
      defaultContest={contest}
      defaultSelected={selected}
      handleSubmit={onUpdate}
    />
  );
}

export default AdminEditContestPage;
