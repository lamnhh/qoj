import React from "react";
import AdminContestForm from "../components/AdminContestForm";
import { useParams } from "react-router-dom";
import useSWR from "swr";
import request from "../helpers/request";
import AdminContest from "../models/AdminContest";
import Contest from "../models/Contest";

interface AdminEditContestPageRouterProps {
  id: string;
}

function AdminEditContestPage() {
  let contestId = parseInt(useParams<AdminEditContestPageRouterProps>().id);
  let { data, error } = useSWR("/api/contest/" + contestId, request);

  if (error) {
    return <section className="error-msg">Server Error</section>;
  }
  if (!data) {
    return <section className="loading-msg">Loading</section>;
  }

  let contest = data as Contest;
  let selected = (data as AdminContest).problemList;

  return (
    <AdminContestForm
      action="Update"
      defaultContest={contest}
      defaultSelected={selected}
      handleSubmit={() => {
        console.log("H");
      }}
    />
  );
}

export default AdminEditContestPage;
