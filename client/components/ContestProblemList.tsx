import React from "react";
import Problem from "../models/Problem";
import ScoreBar from "./ScoreBar";

function ContestProblemList({ problemList }: { problemList: Array<Problem> }) {
  return (
    <table className="problem-list my-table striped">
      <tr>
        <th>#</th>
        <th>Name</th>
        <th>Time limit</th>
        <th>Memory limit</th>
        <th>Your score</th>
      </tr>
      {problemList.map(function(problem, idx) {
        return (
          <tr key={problem.id}>
            <td>{String.fromCharCode("A".charCodeAt(0) + idx)}</td>
            <td>{problem.name}</td>
            <td>
              {problem.timeLimit} second{problem.timeLimit !== 1 ? "s" : ""}
            </td>
            <td>{problem.memoryLimit} MB</td>
            <td>
              <ScoreBar
                maxScore={problem.maxScore}
                testCount={problem.testCount}
              />
            </td>
          </tr>
        );
      })}
    </table>
  );
}

export default ContestProblemList;
