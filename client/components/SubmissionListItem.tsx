import React, { useContext, useEffect, useState, ReactElement } from "react";
import moment from "moment";
import Submission from "../models/Submission";
import WSContext from "../contexts/WSContext";
import WSMessage from "../models/WSMessage";
import ScoreBar from "./ScoreBar";
import { Link } from "react-router-dom";

interface SubmissionListItemProps {
  submission: Submission;
}

function parseSubmissionStatus(status: string): ReactElement {
  let display = status.split("|")[0];
  let tokens = display.split("/").map(a => parseInt(a));
  if (tokens.length === 2) {
    return <ScoreBar maxScore={tokens[0]} testCount={tokens[1]}></ScoreBar>;
  }

  if (display.split(" ")[0] === "Compile") {
    return <span className="status status-ce">{display}</span>;
  }
  return <span className="status status-running">{display}</span>;
}

function SubmissionListItem({ submission }: SubmissionListItemProps) {
  let { socket } = useContext(WSContext);
  let [status, setStatus] = useState<ReactElement>(<></>);

  useEffect(
    function() {
      setStatus(parseSubmissionStatus(submission.status));
    },
    [submission.status]
  );

  useEffect(
    function() {
      function updateStatus(event: MessageEvent) {
        let json: WSMessage = JSON.parse(event.data);
        if (json.submissionId === submission.id) {
          setStatus(parseSubmissionStatus(json.message));
        }
      }

      socket.send(
        JSON.stringify({
          type: "subscribe",
          message: String(submission.id)
        })
      );

      socket.addEventListener("message", updateStatus);
      return function() {
        socket.send(
          JSON.stringify({
            type: "unsubscribe",
            message: String(submission.id)
          })
        );
        socket.removeEventListener("message", updateStatus);
      };
    },
    [submission.id]
  );

  return (
    <tr>
      <td className="id">{submission.id}</td>
      <td className="date">
        {moment(submission.createdAt).format("YYYY-MM-DD hh:mm:ss")}
      </td>
      <td>
        <Link to={"/user/" + submission.username}>{submission.username}</Link>
      </td>
      <td>
        <Link to={"/problem/" + submission.problemId}>
          {submission.problemId} - {submission.problemName}
        </Link>
      </td>
      <td>C++17</td>
      <td className="status-cell">{status}</td>
    </tr>
  );
}

export default SubmissionListItem;
