import React, { useState, useEffect, useContext } from "react";
import request from "../helpers/request";
import ContestItem from "../components/ContestItem";
import Contest from "../models/Contest";
import moment from "moment";
import AppContext from "../contexts/AppContext";

function isPastContest(contest: Contest) {
  let contestEnd = moment(contest.startDate).add(contest.duration, "minutes");
  return contestEnd.isBefore(moment());
}

function ContestListPage() {
  let [contestList, setContestList] = useState<Array<Contest>>([]);
  let { user } = useContext(AppContext);

  useEffect(function() {
    request("/api/c/contest").then(setContestList);
  }, []);

  let upcomingList = contestList.filter(contest => !isPastContest(contest));
  let pastList = contestList.filter(contest => isPastContest(contest));
  return (
    <>
      <header className="page-name align-left-right">
        <h1>Contests</h1>
      </header>
      <section className="contest-list-page align-left-right">
        {upcomingList.length > 0 && (
          <div className="contest-table--wrapper">
            <table className="contest-table my-table full-border">
              <tr className="my-table__title">
                <th colSpan={5}>Running and Upcoming Contests</th>
              </tr>
              <tr className="my-table__header">
                <th className="contest-column">Contest</th>
                <th>Start</th>
                <th>Duration</th>
                <th>Participants</th>
                {!!user && <th className="action"></th>}
              </tr>
              {upcomingList.map(function(contest) {
                return <ContestItem key={contest.id} contest={contest} />;
              })}
            </table>
          </div>
        )}
        <div className="contest-table--wrapper">
          <table className="contest-table my-table full-border">
            <tr className="my-table__title">
              <th colSpan={5}>Past Contests</th>
            </tr>
            <tr className="my-table__header">
              <th className="contest-column">Contest</th>
              <th>Start</th>
              <th>Duration</th>
              <th>Participants</th>
              {!!user && <th className="action"></th>}
            </tr>
            {pastList.map(function(contest) {
              return (
                <ContestItem
                  key={contest.id}
                  contest={contest}
                  showAction={false}
                />
              );
            })}
          </table>
        </div>
      </section>
    </>
  );
}

export default ContestListPage;
