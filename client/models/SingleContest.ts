import Problem from "./Problem";

interface SingleContest {
  id: number;
  name: string;
  problemList: Array<Problem>;
  startDate: Date; // in ISO format, i.e. "2020-02-07T09:44:15.991Z"
  duration: number; // in minutes
}

export default SingleContest;
