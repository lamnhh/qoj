import React from "react";
import useSWR from "swr";
import request from "../helpers/request";
import Contest from "../models/Contest";
import moment from "moment";
import { Link } from "react-router-dom";

function AdminContestPage() {
  let { data } = useSWR("/api/contest", request);

  function deleteContest(contestId: number) {
    if (!confirm("Are you sure you want to discard this contest?")) {
      return;
    }

    request("/api/contest/" + contestId, { method: "DELETE" })
      .then(function() {
        location.reload();
      })
      .catch(function({ error }) {
        alert(error);
      });
  }

  let contestList: Contest[] = data ?? [];
  return (
    <section className="contest-page">
      <h2 className="contest-page__count">Contests: {contestList.length}</h2>
      <table className="admin-table">
        <thead>
          <tr>
            <th>#</th>
            <th>Contest Name</th>
            <th>Start</th>
            <th>Duration</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {contestList.length === 0 && (
            <tr>
              <td colSpan={5}>No item</td>
            </tr>
          )}
          {contestList.map(function(contest) {
            return (
              <tr key={contest.id}>
                <td>{contest.id}</td>
                <td>{contest.name}</td>
                <td>{moment(contest.startDate).format("YYYY-MM-DD hh:mm")}</td>
                <td>
                  {moment
                    .duration(contest.duration, "minutes")
                    .format("HH[:]mm")}
                </td>
                <td className="contest-page__row-action">
                  <Link to={`/contest/edit/${contest.id}`}>
                    <button type="button">Edit</button>
                  </Link>
                  <button
                    type="button"
                    onClick={() => deleteContest(contest.id)}>
                    Discard
                  </button>
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </section>
  );
}

export default AdminContestPage;
