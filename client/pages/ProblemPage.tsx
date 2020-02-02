import React, { useEffect, useState, useRef, useContext } from "react";
import { useParams } from "react-router-dom";
import request from "../helpers/request";
import Problem, { emptyProblem } from "../models/Problem";
import SubmissionList from "../components/SubmissionList";
import AppContext from "../contexts/AppContext";
import ScoreBar from "../components/ScoreBar";
import SubmitForm from "../components/SubmitForm";

interface ProblemPageRouterProps {
  problemId: string;
}

function ProblemPage() {
  let problemId = useParams<ProblemPageRouterProps>().problemId;
  let user = useContext(AppContext).user;
  let isLoggedIn = user !== null;

  let [problem, setProblem] = useState<Problem>(emptyProblem);
  useEffect(
    function() {
      request("/api/problem/" + problemId).then(setProblem);
    },
    [problemId]
  );

  let [tab, setTab] = useState(0);
  let tabListRef = useRef<HTMLDivElement>(null);

  useEffect(
    function() {
      let tabList = tabListRef.current!;
      let tabs = tabList.querySelectorAll('[role="tab"]');
      let len = tabs.length;

      function tabNavigator(e: KeyboardEvent) {
        if (e.keyCode === 39) {
          setTab(function(tab) {
            let newTab = (tab + 1) % len;
            return newTab;
          });
        } else if (e.keyCode === 37) {
          setTab(function(tab) {
            let newTab = (tab + len - 1) % len;
            return newTab;
          });
        }
      }

      tabList.addEventListener("keydown", tabNavigator);
      return function() {
        tabList.removeEventListener("keydown", tabNavigator);
      };
    },
    [isLoggedIn]
  );

  useEffect(
    function() {
      let tabList = tabListRef.current!;
      let tabs = tabList.querySelectorAll('[role="tab"]');
      (tabs[tab] as HTMLElement).focus();
    },
    [isLoggedIn, tab]
  );

  return (
    <>
      <header className="page-name align-left-right">
        <h1>View Problem - {problem.name}</h1>
      </header>
      <section className="problem-page align-left-right">
        <h1 className="problem-page__problem-name">
          {problem.id}. {problem.code} - {problem.name}
        </h1>
        <div
          className="problem-page__tablist"
          ref={tabListRef}
          role="tablist"
          aria-label="Actions">
          <button
            className="problem-page__tab"
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
            <React.Fragment>
              <button
                className="problem-page__tab"
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
              <button
                className="problem-page__tab"
                type="button"
                role="tab"
                aria-label="My submissions"
                aria-selected={tab === 2 ? "true" : "false"}
                aria-controls="tab-mine"
                id="mine"
                onClick={() => setTab(2)}
                tabIndex={tab !== 2 ? -1 : undefined}>
                My Submissions
              </button>
            </React.Fragment>
          )}
          <button
            className="problem-page__tab"
            type="button"
            role="tab"
            aria-label="View all submissions of this problem"
            aria-selected={tab === (isLoggedIn ? 3 : 1) ? "true" : "false"}
            aria-controls="tab-submission"
            id="submission"
            onClick={() => setTab(isLoggedIn ? 3 : 1)}
            tabIndex={tab !== (isLoggedIn ? 3 : 1) ? -1 : undefined}>
            Status
          </button>
        </div>
        <div className="problem-page__tabpanel">
          <div
            tabIndex={0}
            role="tabpanel"
            id="tab-constraints"
            aria-labelledby="constraints"
            hidden={tab !== 0}>
            <table className="my-table problem-page__constraints">
              <tr className="my-table__header">
                <th>Time limit</th>
                <th>Memory limit</th>
                <th className="score">Your score</th>
              </tr>
              <tr>
                <td>
                  {problem.timeLimit} second{problem.timeLimit !== 1 ? "s" : ""}
                </td>
                <td>{problem.memoryLimit} MB</td>
                <td className="score">
                  <ScoreBar
                    maxScore={problem.maxScore}
                    testCount={problem.testCount}
                  />
                </td>
              </tr>
            </table>
          </div>
          {isLoggedIn && (
            <React.Fragment>
              <div
                tabIndex={0}
                role="tabpanel"
                id="tab-submit"
                aria-labelledby="submit"
                hidden={tab !== 1}>
                <SubmitForm problemId={problemId}></SubmitForm>
              </div>
              <div
                tabIndex={0}
                role="tabpanel"
                id="tab-mine"
                aria-labelledby="mine"
                hidden={tab !== 2}>
                {tab === 2 && (
                  <SubmissionList
                    params={[
                      ["problemId", String(problemId)],
                      ["username", user!.username]
                    ]}
                  />
                )}
              </div>
            </React.Fragment>
          )}
          <div
            tabIndex={0}
            role="tabpanel"
            id="tab-submission"
            aria-labelledby="submission"
            hidden={tab !== (isLoggedIn ? 3 : 1)}>
            {tab === (isLoggedIn ? 3 : 1) && (
              <SubmissionList params={[["problemId", String(problemId)]]} />
            )}
          </div>
        </div>
      </section>
    </>
  );
}

export default ProblemPage;
