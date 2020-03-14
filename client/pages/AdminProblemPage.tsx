import React, { useContext, useCallback } from "react";
import useSWR from "swr";
import qs from "querystring";
import request from "../helpers/request";
import Problem from "../models/Problem";
import AppContext from "../contexts/AppContext";
import { Link, useLocation } from "react-router-dom";

function AdminProblemPage() {
  let search = qs.parse(useLocation().search.slice(1)).search ?? "";

  let { user } = useContext(AppContext);
  let { data, error } = useSWR(
    "/api/problem" + (search ? "?search=" + search : ""),
    request
  );
  let problemList: Problem[] = data ?? [];

  let deleteProblem = useCallback(function(problemId: number) {
    if (!confirm("Are you sure you want to discard this problem?")) {
      return;
    }
    request("/api/problem/" + problemId, { method: "DELETE" })
      .then(function() {
        location.reload();
      })
      .catch(function({ error }) {
        alert(error);
      });
  }, []);

  if (error) {
    return <section className="error-msg">Server error</section>;
  }

  return (
    <section className="problem-page">
      <h2 className="problem-page__count">Problems: {problemList.length}</h2>
      <table className="admin-table">
        <thead>
          <tr>
            <th>#</th>
            <th>Name</th>
            <th>Problem Setter</th>
            <th>Time Limit</th>
            <th>Memory Limit</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {problemList.map(function(problem) {
            return (
              <tr key={problem.id}>
                <td>{problem.id}</td>
                <td>
                  {problem.code === problem.name
                    ? problem.code
                    : problem.code + " - " + problem.name}
                </td>
                <td>{user!.username}</td>
                <td>{Math.round(problem.timeLimit * 1000)} ms</td>
                <td>{problem.memoryLimit} MB</td>
                <td className="problem-page__row-action">
                  <Link to={"/problem/edit/" + problem.id}>
                    <button type="button">Edit</button>
                  </Link>
                  <button
                    type="button"
                    onClick={function() {
                      deleteProblem(problem.id);
                    }}>
                    Discard
                  </button>
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </section>
  );
}

export default AdminProblemPage;
