import React, { useCallback } from "react";
import { useParams, useHistory } from "react-router-dom";
import AdminProblemForm from "../components/AdminProblemForm";
import useSWR from "swr";
import request from "../helpers/request";
import Problem from "../models/Problem";
import AdminProblem from "../models/AdminProblem";
import Loading from "../components/Loading";

interface AdminEditProblemPageRouterProps {
  id: string;
}

function AdminEditProblemPage() {
  let problemId = parseInt(useParams<AdminEditProblemPageRouterProps>().id);
  let { data } = useSWR("/api/problem/" + problemId, request);

  let history = useHistory();

  let onUpdate = useCallback(
    function(problem: AdminProblem) {
      let arr = [];

      arr.push(
        request("/api/problem/" + problemId, {
          method: "PATCH",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            code: problem.code,
            name: problem.name,
            timeLimit: problem.timeLimit,
            memoryLimit: problem.memoryLimit
          })
        })
      );

      if (problem.file) {
        let body = new FormData();
        body.append("file", problem.file);
        arr.push(
          request(`/api/problem/${problemId}/test?replace=1`, {
            method: "PUT",
            body
          })
        );
      }

      Promise.all(arr)
        .then(function() {
          history.push("/problem");
        })
        .catch(function({ error }) {
          alert(error);
        });
    },
    [problemId, history]
  );

  if (!data) {
    return <Loading />;
  }

  let currentProblem: Problem = {
    ...data,
    timeLimit: Math.round(data.timeLimit * 1000)
  };
  return (
    <AdminProblemForm
      title={"Edit problem #" + problemId}
      action="Update"
      defaultProblem={currentProblem}
      requireFile={false}
      handleSubmit={onUpdate}
    />
  );
}

export default AdminEditProblemPage;
