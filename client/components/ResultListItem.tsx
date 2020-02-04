import React from "react";
import Result from "../models/Result";

function ResultListItem({ result, index }: { result: Result; index: number }) {
  let verdictColor = result.verdict === "Accepted" ? "green" : "red";
  return (
    <div className="test-result">
      <h3 className="test-result__header">#{index}</h3>
      <div className="test-result__body">
        <p className="test-result__metadata">
          <h4>
            <strong>Time:</strong> {result.executionTime} ms
          </h4>
          .{" "}
          <h4>
            <strong>Memory:</strong> {result.memoryUsed} KB
          </h4>
          .{" "}
          <h4>
            <strong>Score:</strong> {result.score}
          </h4>
          .
        </p>
        <h4 className={`test-result__verdict ${verdictColor}`}>
          Verdict: {result.verdict}
        </h4>
        <h4 className="test-result__preview">Input</h4>
        <pre>{result.inputPreview}</pre>
        <h4 className="test-result__preview">Participant's output</h4>
        <pre>{result.answerPreview}</pre>
        <h4 className="test-result__preview">Jury's output</h4>
        <pre>{result.outputPreview}</pre>
      </div>
    </div>
  );
}

export default ResultListItem;
