import React from "react";
import Submission from "../models/Submission";
import SubmissionListItem from "./SubmissionListItem";
import PaginationProps from "../models/PaginationProps";
import Pagination from "./Pagination";

interface SubmissionListProps {
  socket: WebSocket;
  submissionList: Array<Submission>;
  paginationProps: PaginationProps;
}

function SubmissionList({
  socket,
  paginationProps,
  submissionList
}: SubmissionListProps) {
  return (
    <div className="submission-list--wrapper">
      <table className="submission-list my-table striped">
        <tr className="my-table__header">
          <th className="id">#</th>
          <th className="date">Submission time</th>
          <th>Handle</th>
          <th>Problem</th>
          <th>Language</th>
          <th className="status-cell">Result</th>
          <th>Execution time</th>
          <th>Memory</th>
        </tr>
        {submissionList.map(submission => {
          return (
            <SubmissionListItem
              key={submission.id}
              submission={submission}
              socket={socket}
            />
          );
        })}
      </table>
      <Pagination {...paginationProps} />
    </div>
  );
}

export default SubmissionList;
