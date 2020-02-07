import React, { useCallback } from "react";
import MultipleContest from "../models/MultipleContest";
import moment from "moment";
import { Link, useHistory } from "react-router-dom";
import request from "../helpers/request";

function ContestItem({ contest }: { contest: MultipleContest }) {
  let ms = moment.duration(contest.duration, "minutes").asMilliseconds();
  let duration = moment.utc(ms).format("HH:mm");
  let history = useHistory();

  let onJoin = useCallback(
    function() {
      request(`/api/contest/${contest.id}/register`, { method: "POST" })
        .then(function() {
          history.push("/contest/" + contest.id);
        })
        .catch(function() {
          alert("Something went wrong. Please try again.");
        });
    },
    [contest.id, history]
  );

  return (
    <tr>
      <td className="contest-column">{contest.name}</td>
      <td>{moment(contest.startDate).format("MMM/DD/YYYY, HH:mm")}</td>
      <td>{duration}</td>
      <td>
        <Link to={`/contest/${contest.id}/participants`}>
          <i className="fa fa-user"></i> {contest.numberOfParticipants}
        </Link>
      </td>
      <td>
        {!contest.isRegistered ? (
          <button type="button" className="join-btn" onClick={onJoin}>
            Join Contest
          </button>
        ) : (
          <span className="register-state">Registration completed</span>
        )}
      </td>
    </tr>
  );
}

export default ContestItem;
