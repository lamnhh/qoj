import React from "react";

interface ScoreBarProps {
  maxScore: number;
  testCount: number;
}

function ScoreBar({ maxScore, testCount }: ScoreBarProps) {
  return (
    <div className="problemset-score--wrapper">
      <div className="problemset-score">
        <span>
          {maxScore} / {testCount}
        </span>
        <div className="problemset-score__progress--wrapper">
          <div
            className="problemset-score__progress"
            style={{
              width: `${(maxScore / testCount) * 100}%`
            }}></div>
        </div>
      </div>
    </div>
  );
}

export default ScoreBar;
