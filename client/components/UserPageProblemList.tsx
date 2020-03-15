import React, { useState, useEffect } from "react";
import request from "../helpers/request";
import { Link } from "react-router-dom";

// This interface only stores `id` and `code` of problems
interface PartialProblem {
  id: number;
  code: string;
  isPlaceholder: boolean;
}

function ProblemListItem({ problem }: { problem: PartialProblem }) {
  return (
    <div className="user-page__prob-list__item">
      {!problem.isPlaceholder && (
        <Link to={`/problem/${problem.id}`}>{problem.code}</Link>
      )}
    </div>
  );
}

function UserPageProblemList({ title, url }: { title: string; url: string }) {
  let [problemList, setProblemList] = useState<Array<PartialProblem>>([]);
  useEffect(
    function() {
      request(url)
        .then(function(problemList) {
          return problemList.map(function(
            problem: PartialProblem
          ): PartialProblem {
            return {
              ...problem,
              isPlaceholder: false
            };
          });
        })
        .then(function(problemList: Array<PartialProblem>) {
          while (problemList.length % 4 !== 0) {
            problemList.push({
              id: 0,
              code: "",
              isPlaceholder: true
            });
          }
          return problemList;
        })
        .then(setProblemList)
        .catch(console.log);
    },
    [url]
  );

  return (
    <div className="user-page__prob-list">
      <div className="user-page__prob-list__title">
        <h1>{title}</h1>
      </div>
      <div
        className={
          "user-page__prob-list__content" +
          (problemList.length === 0 ? " empty" : "")
        }>
        {problemList.length === 0 && (
          <div className="user-page__prob-list__item">No item</div>
        )}
        {problemList.map(function(problem: PartialProblem) {
          return <ProblemListItem key={problem.id} problem={problem} />;
        })}
      </div>
    </div>
  );
}

export default UserPageProblemList;
