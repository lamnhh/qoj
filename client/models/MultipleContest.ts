import Problem from "./Problem";

interface MultipleContest {
  id: number;
  name: string;
  startDate: Date; // in ISO format, i.e. "2020-02-07T09:44:15.991Z"
  duration: number; // in minutes
  numberOfParticipants: number;
  isRegistered: boolean;
}

export default MultipleContest;
