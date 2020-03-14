import React, { useState, useEffect } from "react";
import Result from "../models/Result";
import request from "../helpers/request";
import ResultListItem from "./ResultListItem";

function ResultList({ submissionId }: { submissionId: string }) {
  let [resultList, setResultList] = useState<Array<Result>>([]);

  useEffect(
    function() {
      request(`/api/c/submission/${submissionId}/result`).then(setResultList);
    },
    [submissionId]
  );

  return (
    <>
      {resultList.map(function(result, index) {
        return <ResultListItem key={index} result={result} index={index + 1} />;
      })}
    </>
  );
}

export default ResultList;
