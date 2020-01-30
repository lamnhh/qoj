interface Submission {
  id: number;
  username: string;
  problemId: number;
  problemName: string;
  createdAt: Date;
  status: string;
}

export default Submission;
