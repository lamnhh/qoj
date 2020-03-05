import React, { useState, useEffect, useRef } from "react";
import { useParams } from "react-router-dom";
import SubmissionListItem from "../components/SubmissionListItem";
import Submission from "../models/Submission";
import request from "../helpers/request";
import CodeMirror, { EditorFromTextArea } from "codemirror";
import ResultList from "../components/ResultList";

interface SubmissionPageRouterProps {
  submissionId: string;
}

function SubmissionPage() {
  let submissionId = useParams<SubmissionPageRouterProps>().submissionId;
  let [submission, setSubmission] = useState<Submission | null>(null);
  let [code, setCode] = useState("");
  let [compileMessage, setCompileMessage] = useState("");

  let codeRef = useRef<HTMLTextAreaElement | null>(null);
  let editor = useRef<EditorFromTextArea | null>(null);

  useEffect(
    function() {
      request("/api/submission/" + submissionId).then(setSubmission);
      request("/api/submission/" + submissionId + "/code")
        .then(({ code }) => code)
        .then(setCode);
      request("/api/submission/" + submissionId + "/compile")
        .then(({ compileMessage }) => compileMessage)
        .then(setCompileMessage);
    },
    [submissionId]
  );

  useEffect(function() {
    if (codeRef.current) {
      editor.current = CodeMirror.fromTextArea(codeRef.current, {
        lineNumbers: true,
        readOnly: "nocursor",
        mode: "text/x-c++src"
      });
    }
  }, []);

  useEffect(
    function() {
      if (editor.current) {
        editor.current.setValue(code);
      }
    },
    [code]
  );

  return (
    <>
      <header className="page-name align-left-right">
        <h1>Submission #{submissionId}</h1>
      </header>
      <section className="submission-page align-left-right">
        <table className="submission-list my-table">
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
          {submission && <SubmissionListItem submission={submission} />}
        </table>
        <textarea ref={codeRef}></textarea>
        {compileMessage.length > 0 && (
          <div className="submission-page__compile-msg">
            <h2>Compilation message</h2>
            <pre>{compileMessage}</pre>
          </div>
        )}
        <ResultList submissionId={submissionId} />
      </section>
    </>
  );
}

export default SubmissionPage;
