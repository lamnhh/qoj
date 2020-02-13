import React, { useState, useEffect } from "react";
import Problem from "../models/Problem";
import request from "../helpers/request";
import { Link } from "react-router-dom";
import ContestRankingCell from "./ContestRankingCell";

interface ContestRankingProps {
  contestId: string;
  problemList: Array<Problem>;
}

interface Score {
  username: string;
  scoreSum: number;
  scoreList: Array<{ problemId: number; score: number }>;
}

function ContestRanking({ contestId, problemList }: ContestRankingProps) {
  let [scoreList, setScoreList] = useState<Array<Score>>([]);
  useEffect(
    function() {
      request(`/api/contest/${contestId}/score`).then(setScoreList);
    },
    [contestId]
  );

  let maxTotalScore = problemList.reduce(function(acc, { maxScore }) {
    return acc + maxScore;
  }, 0);
  return (
    <div className="contest-ranking">
      <table className="my-table full-border">
        <tr className="my-table__header">
          <th className="index">#</th>
          <th>Username</th>
          {problemList.map(function(problem, idx) {
            return (
              <th key={idx} className="score score-header">
                <span>{String.fromCharCode("A".charCodeAt(0) + idx)}</span>
                <span className="max-score">{problem.maxScore}</span>
              </th>
            );
          })}
          <th className="score score-header">
            <span>Total</span>
            <span className="max-score">{maxTotalScore}</span>
          </th>
        </tr>
        {scoreList.map(function(score, index) {
          let scoreMap: { [id: number]: number } = {};
          score.scoreList.forEach(function({ problemId, score }) {
            scoreMap[problemId] = score;
          });
          return (
            <tr key={score.username}>
              <td className="index">{index + 1}</td>
              <td>
                <Link to={`/user/${score.username}`}>{score.username}</Link>
              </td>
              {problemList.map(function(problem) {
                return (
                  <ContestRankingCell
                    score={scoreMap[problem.id]}
                    maxScore={problem.maxScore}
                  />
                );
              })}
              <td className="score">{score.scoreSum}</td>
            </tr>
          );
        })}
      </table>
    </div>
  );
}

export default ContestRanking;
