import React from "react";
import request from "../helpers/request";
import { useHistory } from "react-router-dom";
import AdminProblemForm from "../components/AdminProblemForm";
import { defaultProblem } from "../models/Problem";
import AdminProblem from "../models/AdminProblem";

function AdminCreateProblemPage() {
  let history = useHistory();
  function onCreate(problem: AdminProblem) {
    let body = new FormData();
    body.append("code", problem.code);
    body.append("name", problem.name);
    body.append("timeLimit", problem.timeLimit.toString());
    body.append("memoryLimit", problem.memoryLimit.toString());
    body.append("file", problem.file!);

    request("/api/problem", {
      method: "POST",
      body
    })
      .then(function() {
        history.push("/problem");
      })
      .catch(function({ error }) {
        alert(error);
      });
  }

  return (
    <AdminProblemForm
      title="Create problem"
      action="Create"
      defaultProblem={defaultProblem}
      requireFile={true}
      handleSubmit={onCreate}
    />
  );
}

export default AdminCreateProblemPage;
