import React from "react";
import useSWR from "swr";
import request from "../helpers/request";
import Contest from "../models/Contest";
import moment from "moment";
import { Link } from "react-router-dom";

function AdminContestPage() {
  let { data, error } = useSWR("/api/contest", request);

  if (error) {
    return <section className="error-msg">Server Error</section>;
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
            <th>Ready</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
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
                <td style={{ color: "red", fontWeight: "bold" }}>Yes</td>
                <td className="contest-page__row-action">
                  <Link to={`/contest/edit/${contest.id}`}>
                    <button type="button">Edit</button>
                  </Link>
                  <button type="button">Discard</button>
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
