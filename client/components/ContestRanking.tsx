import React, { useState, useEffect } from "react";
import Problem from "../models/Problem";
import request from "../helpers/request";
import { Link } from "react-router-dom";
import FlipMove from "react-flip-move";
import ContestRankingCell from "./ContestRankingCell";
import WSContestMessage from "../models/WSContestMessage";

interface ContestRankingProps {
  contestId: string;
  problemList: Array<Problem>;
}

interface Score {
  username: string;
  scoreSum: number;
  scoreList: Array<{ problemId: number; score: number }>;
}

function updateNewSubmission(
  scoreList: Array<Score>,
  submission: WSContestMessage
) {
  let newList;
  if (scoreList.find(score => score.username === submission.username)) {
    newList = scoreList.map(function(score) {
      if (score.username !== submission.username) {
        return score;
      }
      let newScoreList = score.scoreList.map(function(a) {
        if (a.problemId !== submission.problemId) {
          return a;
        }
        return {
          ...a,
          score: Math.max(a.score, submission.score)
        };
      });
      if (
        !newScoreList.find(
          ({ problemId }) => problemId === submission.problemId
        )
      ) {
        newScoreList.push({
          problemId: submission.problemId,
          score: submission.score
        });
      }
      return {
        username: score.username,
        scoreList: newScoreList,
        scoreSum: newScoreList.reduce((acc, cur) => acc + cur.score, 0)
      };
    });
  } else {
    newList = scoreList.concat({
      username: submission.username,
      scoreSum: submission.score,
      scoreList: [
        {
          problemId: submission.problemId,
          score: submission.score
        }
      ]
    });
  }
  return newList.sort(function(a, b) {
    if (a.scoreSum === b.scoreSum) {
      if (a.username < b.username) {
        return -1;
      }
      if (a.username > b.username) {
        return 1;
      }
      return 0;
    }
    return -a.scoreSum + b.scoreSum;
  });
}

function ContestRanking({ contestId, problemList }: ContestRankingProps) {
  let [scoreList, setScoreList] = useState<Array<Score>>([]);
  useEffect(
    function() {
      request(`/api/contest/${contestId}/score`).then(setScoreList);
    },
    [contestId]
  );

  let [socket] = useState(function() {
    return new WebSocket(`ws://${location.host}/ws/contest`);
  });
  let [loading, setLoading] = useState(true);

  useEffect(function() {
    function updateScoreList(event: MessageEvent) {
      let json: WSContestMessage = JSON.parse(event.data);
      setScoreList(function(scoreList) {
        return updateNewSubmission(scoreList, json);
      });
    }

    socket.onopen = function() {
      setLoading(false);
      socket.send(
        JSON.stringify({
          type: "subscribe",
          message: contestId
        })
      );
      socket.addEventListener("message", updateScoreList);
    };

    return function() {
      socket.send(
        JSON.stringify({
          type: "unsubscribe",
          message: contestId
        })
      );
      socket.removeEventListener("message", updateScoreList);
    };
  }, []);

  if (loading) {
    return null;
  }

  let maxTotalScore = problemList.reduce(function(acc, { testCount }) {
    return acc + testCount;
  }, 0);
  return (
    <div className="contest-ranking">
      <table className="my-table full-border">
        <FlipMove className="contest-ranking__flip">
          <tr className="my-table__header">
            <th className="index">#</th>
            <th>Username</th>
            {problemList.map(function(problem, idx) {
              return (
                <th key={idx} className="score score-header">
                  <span>{String.fromCharCode("A".charCodeAt(0) + idx)}</span>
                  <span className="max-score">{problem.testCount}</span>
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
        </FlipMove>
      </table>
    </div>
  );
}

export default ContestRanking;
