import React from "react";
import Submission from "../models/Submission";
import SubmissionListItem from "./SubmissionListItem";
import PaginationProps from "../models/PaginationProps";
import Pagination from "./Pagination";

interface SubmissionListProps {
  title?: string;
  borderStyle?: string;
  socket: WebSocket;
  submissionList: Array<Submission>;
  paginationProps: PaginationProps | null;
}

function SubmissionList({
  title,
  borderStyle = "striped",
  socket,
  paginationProps,
  submissionList
}: SubmissionListProps) {
  return (
    <div className="submission-list--wrapper">
      <table className={`submission-list my-table ${borderStyle}`}>
        {typeof title !== "undefined" && (
          <tr className="my-table__title">
            <th colSpan={8}>{title}</th>
          </tr>
        )}
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
        {submissionList.length === 0 ? (
          <tr className="my-table__empty-row">
            <td colSpan={8}>No item</td>
          </tr>
        ) : (
          submissionList.map(submission => {
            return (
              <SubmissionListItem
                key={submission.id}
                submission={submission}
                socket={socket}
              />
            );
          })
        )}
      </table>
      {paginationProps !== null && <Pagination {...paginationProps} />}
    </div>
  );
}

export default SubmissionList;
