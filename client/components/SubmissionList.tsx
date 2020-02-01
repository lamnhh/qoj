import React from "react";
import Submission from "../models/Submission";
import SubmissionListItem from "./SubmissionListItem";
import WSContext from "../contexts/WSContext";
import request from "../helpers/request";

interface SubmissionListProps {
  baseUrl: string;
}

interface SubmissionListState {
  loading: boolean;
  submissionList: Array<Submission>;
}

class SubmissionList extends React.Component<
  SubmissionListProps,
  SubmissionListState
> {
  socket = new WebSocket("ws://localhost:3000/ws");
  state: SubmissionListState = { loading: true, submissionList: [] };

  fetchSubmissionList = () => {
    request(this.props.baseUrl).then((submissionList: Array<Submission>) => {
      this.setState({ submissionList });
    });
  };

  componentDidMount() {
    this.socket.onopen = () => {
      this.setState({ loading: false });
    };

    request(this.props.baseUrl).then((submissionList: Array<Submission>) => {
      this.setState({ submissionList });
    });
  }

  componentDidUpdate(prevProps: SubmissionListProps) {
    if (prevProps.baseUrl !== this.props.baseUrl) {
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
        <table className="submission-list my-table">
          <tr>
            <th className="id">#</th>
            <th className="date">Submission time</th>
            <th>Handle</th>
            <th>Problem</th>
            <th>Language</th>
            <th className="status-cell">Result</th>
          </tr>
          <WSContext.Provider value={{ socket: this.socket }}>
            {submissionList.map(function(submission) {
              return (
                <SubmissionListItem
                  key={submission.id}
                  submission={submission}></SubmissionListItem>
              );
            })}
          </WSContext.Provider>
        </table>
      </div>
    );
  }
}

export default SubmissionList;
