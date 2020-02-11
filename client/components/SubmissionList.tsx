import React from "react";
import Submission from "../models/Submission";
import SubmissionListItem from "./SubmissionListItem";
import { requestWithHeaders } from "../helpers/request";
import { RouteComponentProps, withRouter } from "react-router-dom";
import qs from "querystring";
import { parsePage, buildURL } from "../helpers/common-helper";
import Pagination from "./Pagination";

interface SubmissionListQuery extends qs.ParsedUrlQuery {
  page: string;
}

interface SubmissionListProps extends RouteComponentProps {
  params: Array<[string, string | string[]]>;
}

interface SubmissionListState {
  page: number;
  submissionCount: number;
  loading: boolean;
  submissionList: Array<Submission>;
}

const pageSize = 15;

class SubmissionList extends React.Component<
  SubmissionListProps,
  SubmissionListState
> {
  socket = new WebSocket("ws://localhost:3000/ws");

  constructor(props: SubmissionListProps) {
    super(props);
    this.state = {
      page: -1,
      submissionCount: 0,
      loading: true,
      submissionList: []
    };
  }

  fetchCurrentPage = () => {
    let queries = qs.parse(
      this.props.location.search.slice(1)
    ) as SubmissionListQuery;
    this.setState({ page: parsePage(queries.page) });
  };

  fetchSubmissionList = () => {
    let { page } = this.state;
    if (page === -1) {
      return;
    }

    let url = buildURL(
      "/api/submission",
      this.props.params.concat([
        ["page", String(page)],
        ["size", String(pageSize)]
      ])
    );
    requestWithHeaders(url).then(
      ([submissionList, headers]: [Array<Submission>, Headers]) => {
        this.setState({
          submissionList,
          submissionCount: parseInt(headers.get("x-count")!)
        });
      }
    );
  };

  componentDidMount() {
    this.socket.onopen = () => {
      this.setState({ loading: false });
    };

    this.fetchCurrentPage();
  }

  componentDidUpdate(
    prevProps: SubmissionListProps,
    prevState: SubmissionListState
  ) {
    if (prevProps.location.search !== this.props.location.search) {
      this.fetchCurrentPage();
    }
    if (prevState.page !== this.state.page) {
      this.fetchSubmissionList();
    }
    if (prevProps.params !== this.props.params) {
      this.fetchSubmissionList();
    }
  }

  componentWillUnmount() {
    this.socket.close();
  }

  render() {
    if (this.state.loading) {
      return null;
    }

    let { submissionList } = this.state;
    return (
      <div className="submission-list--wrapper">
        <table className="submission-list my-table striped">
          <tr className="my-table__header">
            <th className="id">#</th>
            <th className="date">Submission time</th>
            <th>Handle</th>
            <th>Problem</th>
            <th>Language</th>
            <th className="status-cell">Result</th>
            <th>Execution time</th>
            <th>Memory</th>
          </tr>
          {submissionList.map(submission => {
            return (
              <SubmissionListItem
                key={submission.id}
                submission={submission}
                socket={this.socket}
              />
            );
          })}
        </table>
        <Pagination
          totalCount={this.state.submissionCount}
          currentPage={this.state.page}
          pageSize={pageSize}
          onPageChange={(page: number) => {
            this.props.history.push({
              pathname: this.props.location.pathname,
              search: "?page=" + page
            });
          }}
        />
      </div>
    );
  }
}

let SubmissionListWithRouter = withRouter(SubmissionList);

export default SubmissionListWithRouter;
