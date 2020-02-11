import React, { useState, useEffect } from "react";
import request from "../helpers/request";
import ContestItem from "../components/ContestItem";
import Contest from "../models/Contest";
import moment from "moment";

function isPastContest(contest: Contest) {
  let contestEnd = moment(contest.startDate).add(contest.duration, "minutes");
  return contestEnd.isBefore(moment());
}

function ContestListPage() {
  let [contestList, setContestList] = useState<Array<Contest>>([]);

  useEffect(function() {
    request("/api/contest").then(setContestList);
  }, []);

  let upcomingList = contestList.filter(contest => !isPastContest(contest));
  let pastList = contestList.filter(contest => isPastContest(contest));
  return (
    <>
      <header className="page-name align-left-right">
        <h1>Contests</h1>
      </header>
      <section className="contest-list-page align-left-right">
        <table className="contest-table my-table full-border">
          <tr className="my-table__title">
            <th colSpan={5}>Running and Upcoming Contests</th>
          </tr>
          <tr className="my-table__header">
            <th className="contest-column">Contest</th>
            <th>Start</th>
            <th>Duration</th>
            <th>Participants</th>
            <th className="action"></th>
          </tr>
          {upcomingList.map(function(contest) {
            return <ContestItem key={contest.id} contest={contest} />;
          })}
        </table>
        <table className="contest-table my-table full-border">
          <tr className="my-table__title">
            <th colSpan={5}>Past Contests</th>
          </tr>
          <tr className="my-table__header">
            <th className="contest-column">Contest</th>
            <th>Start</th>
            <th>Duration</th>
            <th>Participants</th>
            <th className="action"></th>
          </tr>
          {pastList.map(function(contest) {
            return <ContestItem key={contest.id} contest={contest} />;
          })}
        </table>
      </section>
    </>
  );
}

export default ContestListPage;
