import React from "react";
import useSubmissionList from "../hooks/useSubmissionList";
import SubmissionList from "./SubmissionList";

interface SubmissionListProps {
  params: Array<[string, string | string[]]>;
}

function SubmissionListWrapper({ params }: SubmissionListProps) {
  let { socket, submissionList, paginationProps, loading } = useSubmissionList({
    params
  });
  if (loading) {
    return null;
  }

  return (
    <SubmissionList
      socket={socket}
      submissionList={submissionList}
      paginationProps={paginationProps}
    />
  );
}

export default SubmissionListWrapper;
