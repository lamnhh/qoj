import React, { useState, useCallback, FormEvent, useEffect } from "react";
import Problem from "../models/Problem";
import request from "../helpers/request";
import { useHistory } from "react-router-dom";

interface FormElements extends HTMLFormElement {
  problemId: HTMLSelectElement;
  file: HTMLInputElement;
}

function ProblemsetPage() {
  let [problemList, setProblemList] = useState<Array<Problem>>([]);

  let history = useHistory();

  let handleSubmit = useCallback(function(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    let form = event.target as FormElements;

    let problemId = form.problemId.value;
    let file = form.file.files![0];

    let body = new FormData();
    body.append("problemId", problemId);
    body.append("file", file);

    request("/api/submission", {
      method: "POST",
      body
    }).then(function() {
      history.push("/status");
    });
  }, []);

  useEffect(function() {
    request("/api/problem").then(function(problemList: Array<Problem>) {
      setProblemList(problemList);
    });
  }, []);

  return (
    <form onSubmit={handleSubmit}>
      <select name="problemId">
        {problemList.map(function(problem) {
          return (
            <option key={problem.id} value={problem.id}>
              {problem.code} - {problem.name}
            </option>
          );
        })}
      </select>
      <input type="file" name="file"></input>
      <button type="submit">Submit</button>
    </form>
  );
}

export default ProblemsetPage;
