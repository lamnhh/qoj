import React, { useState, useEffect } from "react";
import Problem from "../models/Problem";
import request from "../helpers/request";
import { Link } from "react-router-dom";
import ScoreBar from "../components/ScoreBar";

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
      <section className="align-left-right">
        <div className="problemset">
          <table>
            <tr>
              <th>#</th>
              <th>Your score</th>
              <th>Problem code</th>
              <th>Problem title</th>
            </tr>
            {problemList.map(function(problem) {
              return (
                <tr key={problem.id}>
                  <td>{problem.id}</td>
                  <td>
                    <ScoreBar problem={problem}></ScoreBar>
                  </td>
                  <td>
                    <Link to={"/problem/" + problem.id}>{problem.code}</Link>
                  </td>
                  <td>
                    <Link to={"/problem/" + problem.id}>{problem.name}</Link>
                  </td>
                </tr>
              );
            })}
          </table>
        </div>
      </section>
    </>
  );
}

export default ProblemsetPage;
