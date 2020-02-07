import React, { useState, useEffect } from "react";
import request from "../helpers/request";
import ContestItem from "../components/ContestItem";
import Contest from "../models/Contest";

function ContestListPage() {
  let [contestList, setContestList] = useState<Array<Contest>>([]);

  useEffect(function() {
    request("/api/contest").then(setContestList);
  }, []);

  return (
    <>
      <header className="page-name align-left-right">
        <h1>Contests</h1>
      </header>
      <section className="contest-list-page align-left-right">
        <table className="contest-table my-table full-border">
          <tr className="my-table__header">
            <th className="contest-column">Contest</th>
            <th>Start</th>
            <th>Duration</th>
            <th>Participants</th>
            <th></th>
          </tr>
          {contestList.map(function(contest) {
            return <ContestItem key={contest.id} contest={contest} />;
          })}
        </table>
      </section>
    </>
  );
}

export default ContestListPage;
