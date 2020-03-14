import React, { useEffect, useState, useContext, useCallback } from "react";
import { useParams, RouteComponentProps, useHistory } from "react-router-dom";
import request from "../helpers/request";
import Problem, { emptyProblem } from "../models/Problem";
import SubmissionList from "../components/SubmissionListWrapper";
import AppContext from "../contexts/AppContext";
import ScoreBar from "../components/ScoreBar";
import SubmitForm from "../components/SubmitForm";
import { Tabs, TabList, TabPanel, Tab } from "react-tabs";
import SubmissionListWrapper from "../components/SubmissionListWrapper";

interface ProblemPageProps extends RouteComponentProps {
  tab: number;
}

interface ProblemPageRouterProps {
  problemId: string;
}

function ProblemPage({ tab }: ProblemPageProps) {
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

  let history = useHistory();
  let onTabChange = useCallback(
    function(tab) {
      let url = "/problem/" + problemId;
      switch (tab) {
        case 0:
          url += "";
          break;
        case 1:
          url += "/submit";
          break;
        case 2:
          url += "/my";
          break;
        case 3:
          url += "/status";
          break;
      }
      history.push(url);
    },
    [problemId, history]
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
        <Tabs selectedIndex={tab} onSelect={onTabChange}>
          <TabList className="my-tablist">
            <Tab className="my-tab" aria-label="View problem's contraints">
              Contraints
            </Tab>
            <Tab
              className="my-tab"
              aria-label="Submit your solution"
              disabled={!isLoggedIn}>
              Submit
            </Tab>
            <Tab
              className="my-tab"
              aria-label="View your submissions"
              disabled={!isLoggedIn}>
              My Submissions
            </Tab>
            <Tab
              className="my-tab"
              aria-label="View all submissions of this problem">
              Status
            </Tab>
          </TabList>
          <TabPanel>
            <div>
              <table className="my-table problem-page__constraints horizontal">
                <tr className="my-table__header">
                  <th>Time limit</th>
                  <th>Memory limit</th>
                  <th className="score">Your score</th>
                </tr>
                <tr>
                  <td>
                    {problem.timeLimit} second
                    {problem.timeLimit !== 1 ? "s" : ""}
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
              <table className="my-table problem-page__constraints vertical">
                <tr className="my-table__header">
                  <th>Time limit</th>
                  <th>Memory limit</th>
                </tr>
                <tr>
                  <td>
                    {problem.timeLimit} second
                    {problem.timeLimit !== 1 ? "s" : ""}
                  </td>
                  <td>{problem.memoryLimit} MB</td>
                </tr>
                <tr className="my-table__header">
                  <th className="score" colSpan={2}>
                    Your score
                  </th>
                </tr>
                <tr>
                  <td className="score" colSpan={2}>
                    <ScoreBar
                      maxScore={problem.maxScore}
                      testCount={problem.testCount}
                    />
                  </td>
                </tr>
              </table>
            </div>
          </TabPanel>
          <TabPanel>
            <SubmitForm problemList={[problem]} redirectUrl="/status" />
          </TabPanel>
          <TabPanel>
            {user && (
              <SubmissionList
                params={[
                  ["problemId", String(problemId)],
                  ["username", user.username]
                ]}
              />
            )}
          </TabPanel>
          <TabPanel>
            <SubmissionListWrapper
              params={[["problemId", String(problemId)]]}
            />
          </TabPanel>
        </Tabs>
      </section>
    </>
  );
}

export default ProblemPage;
