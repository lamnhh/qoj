import React from "react";
import SubmissionList from "../components/SubmissionList";

function SubmissionPage() {
  return (
    <>
      <header className="page-name align-left-right">
        <h1>Submission Status</h1>
      </header>
      <section className="align-left-right submission-page">
        <SubmissionList baseUrl="/api/submission"></SubmissionList>;
      </section>
    </>
  );
}

export default SubmissionPage;
