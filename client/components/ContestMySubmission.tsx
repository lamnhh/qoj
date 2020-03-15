import React, { useMemo } from "react";
import Problem from "../models/Problem";
import useSubmissionList from "../hooks/useSubmissionList";
import Contest from "../models/Contest";
import Submission from "../models/Submission";
import moment from "moment";
import SubmissionList from "./SubmissionList";

function isInContest(contest: Contest, submission: Submission) {
  let contestEnd = moment(contest.startDate).add(contest.duration, "minutes");
  return contestEnd.isAfter(moment(submission.createdAt));
}

function ContestMySubmission({
  problemList,
  username,
  contest
}: {
  problemList: Array<Problem>;
  username: string;
  contest: Contest;
}) {
  let params = useMemo(
    function(): Array<[string, string | string[]]> {
      return [
        ["problemId", problemList.map(({ id }) => String(id))],
        ["username", username]
      ];
    },
    [problemList, username]
  );

  let { socket, submissionList, paginationProps, loading } = useSubmissionList({
    params,
    pageSize: -1
  });

  if (loading) {
    return null;
  }

  let lateSubmissionList = submissionList.filter(
    submission => !isInContest(contest, submission)
  );
  let inContestSubmissionList = submissionList.filter(submission =>
    isInContest(contest, submission)
  );
  return (
    <div className="my-submission">
      {lateSubmissionList.length > 0 && (
        <SubmissionList
          title="Late submissions"
          borderStyle="full-border"
          socket={socket}
          submissionList={lateSubmissionList}
          paginationProps={paginationProps}
        />
      )}
      <SubmissionList
        title="My in-contest submissions"
        borderStyle="full-border"
        socket={socket}
        submissionList={inContestSubmissionList}
        paginationProps={paginationProps}
      />
    </div>
  );
}

export default ContestMySubmission;
