import React, { useState, useCallback, FormEvent, useEffect } from "react";
import Problem from "../models/Problem";
import request from "../helpers/request";
import { Link } from "react-router-dom";

function ProblemsetPage() {
  let [problemList, setProblemList] = useState<Array<Problem>>([]);
  useEffect(function() {
    request("/api/problem").then(setProblemList);
  }, []);

  return (
    <>
      <header className="page-name align-left-right">
        <h1>Problemset</h1>
      </header>
      <section>
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Score</th>
              <th>Problem title</th>
            </tr>
          </thead>
          <tbody>
            {problemList.map(function(problem) {
              return (
                <tr key={problem.id}>
                  <td>{problem.id}</td>
                  <td>0/100</td>
                  <td>
                    <Link to={"/problem/" + problem.id}>
                      {problem.code} - {problem.name}
                    </Link>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </section>
    </>
  );
}

export default ProblemsetPage;
