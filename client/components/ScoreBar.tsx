import React from "react";
import Problem from "../models/Problem";

function ScoreBar({ problem }: { problem: Problem }) {
  return (
    <div className="problemset-score--wrapper">
      <div className="problemset-score">
        <span>
          {problem.maxScore} / {problem.testCount}
        </span>
        <div className="problemset-score__progress--wrapper">
          <div
            className="problemset-score__progress"
            style={{
              width: `${(problem.maxScore / problem.testCount) * 100}%`
            }}></div>
        </div>
      </div>
    </div>
  );
}

export default ScoreBar;
