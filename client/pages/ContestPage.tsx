import React, { useState, useEffect, useContext, useCallback } from "react";
import { useParams, RouteComponentProps, useHistory } from "react-router-dom";
import request from "../helpers/request";
import Contest from "../models/Contest";
import ContestCountDown from "../components/ContestCountDown";
import { Tabs, TabList, Tab, TabPanel } from "react-tabs";
import moment from "moment";
import Problem from "../models/Problem";
import ContestProblemList from "../components/ContestProblemList";
import SubmitForm from "../components/SubmitForm";
import AppContext from "../contexts/AppContext";
import SubmissionListWrapper from "../components/SubmissionListWrapper";
import ContestMySubmission from "../components/ContestMySubmission";

interface ContestPageProps extends RouteComponentProps {
  tab: number;
}

interface ContestPageRouterProps {
  contestId: string;
}

function ContestPage({ tab }: ContestPageProps) {
  let contestId = useParams<ContestPageRouterProps>().contestId;

  let [contest, setContest] = useState<Contest | null>(null);
  useEffect(
    function() {
      request(`/api/contest/${contestId}`).then(setContest);
    },
    [contestId]
  );

  let [problemList, setProblemList] = useState<Array<Problem> | null>(null);
  useEffect(
    function() {
      request(`/api/contest/${contestId}/problem`).then(setProblemList);
    },
    [contestId]
  );

  let user = useContext(AppContext).user;
  let history = useHistory();
  let onTabChange = useCallback(
    function(newTab) {
      if (tab === newTab) {
        return;
      }
      let url = "/contest/" + contestId;
      switch (newTab) {
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
        case 4:
          url += "/ranking";
          break;
      }
      history.push(url);
    },
    [contestId, history, tab]
  );

  if (contest === null || problemList === null) {
    return null;
  }

  let hasStarted = moment().isAfter(moment(contest.startDate));
  let isLoggedIn = user !== null;
  return (
    <>
      <header className="page-name align-left-right">
        <h1>{contest?.name}</h1>
      </header>
      {!hasStarted ? (
        <ContestCountDown contest={contest}></ContestCountDown>
      ) : (
        <Tabs
          className="contest-page align-left-right"
          selectedIndex={tab}
          onSelect={onTabChange}>
          <TabList className="my-tablist">
            <Tab className="my-tab">Problems</Tab>
            <Tab className="my-tab" disabled={!isLoggedIn}>
              Submit
            </Tab>
            <Tab className="my-tab" disabled={!isLoggedIn}>
              My submissions
            </Tab>
            <Tab className="my-tab">Status</Tab>
            <Tab className="my-tab">Ranking</Tab>
          </TabList>
          <TabPanel>
            <ContestProblemList problemList={problemList}></ContestProblemList>
          </TabPanel>
          <TabPanel>
            <SubmitForm
              problemList={problemList}
              redirectUrl={`/contest/${contestId}/my`}
            />
          </TabPanel>
          <TabPanel>
            {user && (
              <ContestMySubmission
                problemList={problemList}
                username={user.username}
                contest={contest}
              />
            )}
          </TabPanel>
          <TabPanel>
            <SubmissionListWrapper
              params={[["problemId", problemList.map(({ id }) => String(id))]]}
            />
          </TabPanel>
          <TabPanel>Ranking</TabPanel>
        </Tabs>
      )}
    </>
  );
}

export default ContestPage;
