import React from "react";
import Submission from "../models/Submission";
import SubmissionRow from "../components/SubmissionRow";
import WSContext from "../WSContext";
import request from "../helpers/request";

interface SubmissionPageState {
  loading: boolean;
  submissionList: Array<Submission>;
}

class SubmissionPage extends React.Component<{}, SubmissionPageState> {
  socket = new WebSocket("ws://localhost:3000/ws");
  state: SubmissionPageState = { loading: true, submissionList: [] };

  componentDidMount() {
    this.socket.onopen = () => {
      this.setState({ loading: false });
    };

    request("/api/submission").then((submissionList: Array<Submission>) => {
      this.setState({ submissionList });
    });
  }

  render() {
    if (this.state.loading) {
      return null;
    }

    let { submissionList } = this.state;
    return (
      <WSContext.Provider value={{ socket: this.socket }}>
        {submissionList.map(function(submission) {
          return <SubmissionRow key={submission.id} submission={submission}></SubmissionRow>;
        })}
      </WSContext.Provider>
    );
  }
}

export default SubmissionPage;
