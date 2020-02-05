interface Submission {
  id: number;
  username: string;
  problemId: number;
  problemName: string;
  createdAt: Date;
  status: string;
  executionTime: number;
  memoryUsed: number;
  languageId: number;
  language: string;
}

export default Submission;
