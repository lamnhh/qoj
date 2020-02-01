interface Problem {
  id: number;
  code: string;
  name: string;
  timeLimit: number;
  memoryLimit: number;
  maxScore: number;
  testCount: number;
}

const emptyProblem: Problem = {
  id: 0,
  code: "",
  name: "",
  timeLimit: 0,
  memoryLimit: 0,
  maxScore: 0,
  testCount: 0
};

export default Problem;
export { emptyProblem };
