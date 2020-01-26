import React, { useContext, useEffect } from "react";
import moment from "moment";
import Submission from "../models/Submission";
import WSContext from "../WSContext";
import WSMessage from "../models/WSMessage";

interface SubmissionRowProps {
  submission: Submission;
}

let SubmissionRow: React.FC<SubmissionRowProps> = ({ submission }) => {
  let { socket } = useContext(WSContext);
  let [status, setStatus] = React.useState("finished");

  useEffect(
    function() {
      function updateStatus(event: MessageEvent) {
        let json: WSMessage = JSON.parse(event.data);
        if (json.submissionId === submission.id) {
          setStatus(json.message);
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
    <div style={{ display: "grid", gridTemplateColumns: "repeat(5, 1fr)" }}>
      <h3>{submission.id}</h3>
      <h3>{moment(submission.createdAt).format("MMM/DD/YYYY hh:mm:ss")}</h3>
      <h3>{submission.username}</h3>
      <h3>
        {submission.problemId} - {submission.problemName}
      </h3>
      <h3>{status}</h3>
    </div>
  );
};

export default SubmissionRow;
