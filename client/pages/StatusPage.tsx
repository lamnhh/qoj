import React from "react";
import SubmissionListWrapper from "../components/SubmissionListWrapper";

function StatusPage() {
  return (
    <>
      <header className="page-name align-left-right">
        <h1>Submission Status</h1>
      </header>
      <section className="align-left-right status-page">
        <SubmissionListWrapper params={[["allowInContest", "false"]]} />
      </section>
    </>
  );
}

export default StatusPage;
