import React from "react";
import SubmissionList from "../components/SubmissionList";

function StatusPage() {
  return (
    <>
      <header className="page-name align-left-right">
        <h1>Submission Status</h1>
      </header>
      <section className="align-left-right status-page">
        <SubmissionList params={[]}></SubmissionList>
      </section>
    </>
  );
}

export default StatusPage;
