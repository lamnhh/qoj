import React from "react";
import request from "../helpers/request";
import { useHistory } from "react-router-dom";
import AdminProblemForm from "../components/AdminProblemForm";
import { defaultProblem } from "../models/Problem";

function AdminCreateProblemPage() {
  let history = useHistory();
  function onCreate(body: FormData) {
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
      handleSubmit={onCreate}
    />
  );
}

export default AdminCreateProblemPage;
