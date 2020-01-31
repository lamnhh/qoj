import React, {
  useCallback,
  FormEvent,
  useEffect,
  useState,
  useRef,
  useContext
} from "react";
import { useParams, useHistory } from "react-router-dom";
import request from "../helpers/request";
import Problem, { emptyProblem } from "../models/Problem";
import SubmissionList from "../components/SubmissionList";
import AppContext from "../contexts/AppContext";

interface FormElements extends HTMLFormElement {
  file: HTMLInputElement;
}

interface ProblemPageRouterProps {
  problemId: string;
}

function ProblemPage() {
  let history = useHistory();
  let problemId = useParams<ProblemPageRouterProps>().problemId;
  let isLoggedIn = useContext(AppContext).user !== null;

  let [problem, setProblem] = useState<Problem>(emptyProblem);
  useEffect(
    function() {
      request("/api/problem/" + problemId).then(setProblem);
    },
    [problemId]
  );

  let handleSubmit = useCallback(
    function(event: FormEvent<HTMLFormElement>) {
      event.preventDefault();
      let form = event.target as FormElements;
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
    },
    [problemId]
  );

  let [tab, setTab] = useState(0);
  let tabListRef = useRef<HTMLDivElement>(null);

  useEffect(function() {
    let tabList = tabListRef.current!;
    let tabs = tabList.querySelectorAll('[role="tab"]');
    tabList.addEventListener("keydown", function(e) {
      if (e.keyCode === 39) {
        setTab(function(tab) {
          let newTab = (tab + 1) % 3;
          (tabs[newTab] as HTMLElement).focus();
          return newTab;
        });
      } else if (e.keyCode === 37) {
        setTab(function(tab) {
          let newTab = (tab + 2) % 3;
          (tabs[newTab] as HTMLElement).focus();
          return newTab;
        });
      }
    });
  }, []);

  return (
    <section>
      <div ref={tabListRef} role="tablist" aria-label="Actions">
        <button
          type="button"
          role="tab"
          aria-label="View contraints"
          aria-selected={tab === 0 ? "true" : "false"}
          aria-controls="tab-contraints"
          id="constraints"
          onClick={() => setTab(0)}
          tabIndex={tab !== 0 ? -1 : undefined}>
          Contraints
        </button>
        {isLoggedIn && (
          <button
            type="button"
            role="tab"
            aria-label="Submit"
            aria-selected={tab === 1 ? "true" : "false"}
            aria-controls="tab-submit"
            id="submit"
            onClick={() => setTab(1)}
            tabIndex={tab !== 1 ? -1 : undefined}>
            Submit
          </button>
        )}
        <button
          type="button"
          role="tab"
          aria-label="View all submissions of this problem"
          aria-selected={tab === 2 ? "true" : "false"}
          aria-controls="tab-submission"
          id="submission"
          onClick={() => setTab(2)}
          tabIndex={tab !== 2 ? -1 : undefined}>
          Submission
        </button>
      </div>
      <div
        tabIndex={0}
        role="tabpanel"
        id="tab-constraints"
        aria-labelledby="constraints"
        hidden={tab !== 0}>
        <table>
          <tr>
            <th>Time Limit</th>
            <td>
              {problem.timeLimit} second{problem.timeLimit !== 1 ? "s" : ""}
            </td>
          </tr>
          <tr>
            <th>Memory Limit</th>
            <td>{problem.memoryLimit} MB</td>
          </tr>
        </table>
      </div>
      {isLoggedIn && (
        <div
          tabIndex={0}
          role="tabpanel"
          id="tab-submit"
          aria-labelledby="submit"
          hidden={tab !== 1}>
          <form onSubmit={handleSubmit}>
            <input type="file" name="file" required />
            <button type="submit">Submit</button>
          </form>
        </div>
      )}
      <div
        tabIndex={0}
        role="tabpanel"
        id="tab-submission"
        aria-labelledby="submission"
        hidden={tab !== 2}>
        {tab === 2 && (
          <SubmissionList baseUrl={"/api/submission?problemId=" + problemId} />
        )}
      </div>
    </section>
  );
}

export default ProblemPage;
