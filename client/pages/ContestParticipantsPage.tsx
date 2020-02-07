import React, { useState, useEffect } from "react";
import { useParams, Link } from "react-router-dom";
import request from "../helpers/request";
import Contest from "../models/Contest";

interface ContestParticipantsPageRouterProps {
  contestId: string;
}

function ContestParticipantsPage() {
  let contestId = useParams<ContestParticipantsPageRouterProps>().contestId;
  let [contest, setContest] = useState<Contest | null>(null);
  useEffect(
    function() {
      request(`/api/contest/${contestId}`).then(setContest);
    },
    [contestId]
  );

  let [userList, setUserList] = useState<Array<string>>([]);
  useEffect(
    function() {
      request(`/api/contest/${contestId}/participant`).then(setUserList);
    },
    [contestId]
  );

  return (
    <>
      <header className="page-name align-left-right">
        <h1>Registrants for {contest?.name}</h1>
      </header>
      <section className="contest-list-page align-left-right">
        <table className="participant-table my-table full-border">
          <tr className="my-table__header">
            <th>#</th>
            <th>Username</th>
          </tr>
          {userList.map(function(user, idx) {
            return (
              <tr key={user}>
                <td>{idx + 1}</td>
                <td>
                  <Link to={`/user/${user}`}>{user}</Link>
                </td>
              </tr>
            );
          })}
        </table>
      </section>
    </>
  );
}

export default ContestParticipantsPage;
