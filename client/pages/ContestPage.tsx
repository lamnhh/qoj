import React, { useState, useEffect, useContext } from "react";
import { useParams } from "react-router-dom";
import request from "../helpers/request";
import Contest from "../models/Contest";
import ContestCountDown from "../components/ContestCountDown";
import { Tabs, TabList, Tab, TabPanel } from "react-tabs";
import moment from "moment";
import Problem from "../models/Problem";
import ContestProblemList from "../components/ContestProblemList";
import SubmitForm from "../components/SubmitForm";
import SubmissionList from "../components/SubmissionList";
import AppContext from "../contexts/AppContext";

interface ContestPageRouterProps {
  contestId: string;
}

function ContestPage() {
  let contestId = useParams<ContestPageRouterProps>().contestId;

  let [contest, setContest] = useState<Contest | null>(null);
  useEffect(
    function() {
      request(`/api/contest/${contestId}`).then(setContest);
    },
    [contestId]
  );

  let [problemList, setProblemList] = useState<Array<Problem>>([]);
  useEffect(
    function() {
      request(`/api/contest/${contestId}/problem`).then(setProblemList);
    },
    [contestId]
  );

  let user = useContext(AppContext).user;

  if (contest === null) {
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
        <Tabs className="contest-page align-left-right">
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
            <SubmitForm problemList={problemList}></SubmitForm>
          </TabPanel>
          <TabPanel>
            {user && (
              <SubmissionList
                params={[
                  ["problemId", problemList.map(({ id }) => String(id))],
                  ["username", user.username]
                ]}
              />
            )}
          </TabPanel>
          <TabPanel>
            <SubmissionList
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
