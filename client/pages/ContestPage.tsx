import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import request from "../helpers/request";
import Contest from "../models/Contest";

interface ContestPageRouterProps {
  contestId: string;
}

function ContestPage() {
  let contestId = useParams<ContestPageRouterProps>().contestId;
  let [contest, setContest] = useState<Contest | null>(null);
  useEffect(
    function() {
      request(`/api/contest/${contestId}`).then(setContest);
    },
    [contestId]
  );

  return (
    <>
      <header className="page-name align-left-right">
        <h1>Contest {contest?.name}</h1>
      </header>
    </>
  );
}

export default ContestPage;
