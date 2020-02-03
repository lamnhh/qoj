import React, { useEffect, useState, useContext } from "react";
import { useParams } from "react-router-dom";
import request from "../helpers/request";
import Problem, { emptyProblem } from "../models/Problem";
import SubmissionList from "../components/SubmissionList";
import AppContext from "../contexts/AppContext";
import ScoreBar from "../components/ScoreBar";
import SubmitForm from "../components/SubmitForm";
import { Tabs, TabList, TabPanel, Tab } from "react-tabs";

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

  return (
    <>
      <header className="page-name align-left-right">
        <h1>View Problem - {problem.name}</h1>
      </header>
      <section className="problem-page align-left-right">
        <h1 className="problem-page__problem-name">
          {problem.id}. {problem.code} - {problem.name}
        </h1>
        <Tabs>
          <TabList className="my-tablist">
            <Tab className="my-tab" aria-label="View problem's contraints">
              Contraints
            </Tab>
            {isLoggedIn && (
              <Tab className="my-tab" aria-label="Submit your solution">
                Submit
              </Tab>
            )}
            {isLoggedIn && (
              <Tab className="my-tab" aria-label="View your submissions">
                My Submissions
              </Tab>
            )}
            <Tab
              className="my-tab"
              aria-label="View all submissions of this problem">
              Status
            </Tab>
          </TabList>
          <TabPanel>
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
          </TabPanel>
          {isLoggedIn && (
            <TabPanel>
              <SubmitForm problemId={problemId}></SubmitForm>
            </TabPanel>
          )}
          {isLoggedIn && (
            <TabPanel>
              <SubmissionList
                params={[
                  ["problemId", String(problemId)],
                  ["username", user!.username]
                ]}
              />
            </TabPanel>
          )}
          <TabPanel>
            <SubmissionList params={[["problemId", String(problemId)]]} />
          </TabPanel>
        </Tabs>
      </section>
    </>
  );
}

export default ProblemPage;
