import React, { useState, useEffect } from "react";
import Problem from "../models/Problem";
import { requestWithHeaders } from "../helpers/request";
import { Link, useLocation, useHistory } from "react-router-dom";
import ScoreBar from "../components/ScoreBar";
import qs from "querystring";
import Pagination from "../components/Pagination";
import { parsePage } from "../helpers/common-helper";

interface ProblemsetPageQuery extends qs.ParsedUrlQuery {
  page: string;
}

function ProblemsetPage() {
  let [problemCount, setProblemCount] = useState(0);
  let [problemList, setProblemList] = useState<Array<Problem>>([]);

  const pageSize = 20;
  let history = useHistory();
  let queries = qs.parse(useLocation().search.slice(1)) as ProblemsetPageQuery;
  let currentPage = parsePage(queries.page);
  useEffect(
    function() {
      let url = `/api/problem?page=${currentPage}&size=${pageSize}`;
      requestWithHeaders(url).then(function([problemList, headers]) {
        setProblemList(problemList);
        setProblemCount(parseInt(headers.get("x-count") ?? "0"));
      });
    },
    [currentPage]
  );

  return (
    <>
      <header className="page-name align-left-right">
        <h1>Problemset</h1>
      </header>
      <section className="align-left-right">
        <div className="problemset">
          <table className="my-table striped">
            <tr className="my-table__header">
              <th>#</th>
              <th>Your score</th>
              <th>Problem code</th>
              <th>Problem title</th>
            </tr>
            {problemList.map(function(problem) {
              return (
                <tr key={problem.id}>
                  <td>{problem.id}</td>
                  <td>
                    <ScoreBar
                      maxScore={problem.maxScore}
                      testCount={problem.testCount}
                    />
                  </td>
                  <td>
                    <Link to={"/problem/" + problem.id}>{problem.code}</Link>
                  </td>
                  <td>
                    <Link to={"/problem/" + problem.id}>{problem.name}</Link>
                  </td>
                </tr>
              );
            })}
          </table>
          <Pagination
            totalCount={problemCount}
            pageSize={pageSize}
            currentPage={currentPage}
            onPageChange={function(page: number) {
              history.push("/?page=" + page);
            }}
          />
        </div>
      </section>
    </>
  );
}

export default ProblemsetPage;
