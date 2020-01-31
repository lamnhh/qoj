interface Problem {
  id: number;
  code: string;
  name: string;
  timeLimit: number;
  memoryLimit: number;
}

const emptyProblem: Problem = {
  id: 0,
  code: "",
  name: "",
  timeLimit: 0,
  memoryLimit: 0
};

export default Problem;
export { emptyProblem };
