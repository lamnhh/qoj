import { useState, useEffect, useCallback } from "react";
import { buildURL, parsePage } from "../helpers/common-helper";
import { useLocation, useHistory } from "react-router-dom";
import qs from "querystring";
import Submission from "../models/Submission";
import { requestWithHeaders } from "../helpers/request";

interface HookProps {
  params: Array<[string, string | string[]]>;
}

interface SubmissionListQuery extends qs.ParsedUrlQuery {
  page: string;
}

const PAGE_SIZE = 15;

function useSubmissionList({ params }: HookProps) {
  let [socket] = useState<WebSocket>(new WebSocket("ws://localhost:3000/ws"));
  let [loading, setLoading] = useState(true);
  useEffect(function() {
    socket.onopen = function() {
      setLoading(false);
    };
    return function() {
      socket.close();
    };
  }, []);

  let [page, setPage] = useState(1);
  let { pathname, search } = useLocation();
  useEffect(
    function() {
      let queries = qs.parse(search.slice(1)) as SubmissionListQuery;
      setPage(parsePage(queries.page));
    },
    [search]
  );

  let [submissionList, setSubmissionList] = useState<Array<Submission>>([]);
  let [submissionCount, setSubmissionCount] = useState(0);
  useEffect(
    function() {
      if (page === -1) {
        return;
      }

      let url = buildURL(
        "/api/submission",
        params.concat([
          ["page", String(page)],
          ["size", String(PAGE_SIZE)]
        ])
      );
      requestWithHeaders(url).then(
        ([submissionList, headers]: [Array<Submission>, Headers]) => {
          setSubmissionList(submissionList);
          setSubmissionCount(parseInt(headers.get("x-count")!));
        }
      );
    },
    [page, params]
  );

  let history = useHistory();
  let onPageChange = useCallback(
    function(page: number) {
      history.push({
        pathname,
        search: "?page=" + page
      });
    },
    [history, pathname]
  );

  return {
    submissionList,
    paginationProps: {
      onPageChange,
      currentPage: page,
      totalCount: submissionCount,
      pageSize: PAGE_SIZE
    },
    socket,
    loading
  };
}

export default useSubmissionList;
