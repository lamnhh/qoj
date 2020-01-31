import React from "react";
import SubmissionList from "../components/SubmissionList";

function SubmissionPage() {
  return (
    <>
      <header>
        <h1>Submission Status</h1>
      </header>
      <SubmissionList baseUrl="/api/submission"></SubmissionList>;
    </>
  );
}

export default SubmissionPage;
