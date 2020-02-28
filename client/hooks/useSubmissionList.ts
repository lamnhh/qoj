import { useState, useEffect, useCallback } from "react";
import { buildURL, parsePage } from "../helpers/common-helper";
import { useLocation, useHistory } from "react-router-dom";
import qs from "querystring";
import Submission from "../models/Submission";
import { requestWithHeaders } from "../helpers/request";
import PaginationProps from "../models/PaginationProps";

interface HookProps {
  params: Array<[string, string | string[]]>;
  pageSize?: number;
}

interface HookResult {
  submissionList: Array<Submission>;
  paginationProps: PaginationProps | null;
  socket: WebSocket;
  loading: boolean;
}

interface SubmissionListQuery extends qs.ParsedUrlQuery {
  page: string;
}

/**
 * Use pageSize = -1 to fetch all submissions
 * @param param0
 */
function useSubmissionList({ params, pageSize = 15 }: HookProps): HookResult {
  let [socket] = useState<WebSocket>(function() {
    return new WebSocket(`ws://${location.host}/ws/status`);
  });
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
          ["size", String(pageSize)]
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
    paginationProps:
      pageSize === -1
        ? null
        : {
            onPageChange,
            currentPage: page,
            totalCount: submissionCount,
            pageSize
          },
    socket,
    loading
  };
}

export default useSubmissionList;
