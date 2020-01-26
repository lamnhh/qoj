interface Submission {
  id: number;
  username: string;
  problemId: number;
  problemName: string;
  createdAt: Date;
  score: number;
}

export default Submission;
