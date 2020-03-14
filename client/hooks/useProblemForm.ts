import Problem, { defaultProblem } from "../models/Problem";
import { useReducer, useCallback, ChangeEvent } from "react";

type ActionName = "code" | "name" | "timeLimit" | "memoryLimit";
type Action = {
  name: ActionName;
  value: string;
};

function reducer(state: Problem, { name, value }: Action): Problem {
  switch (name) {
    case "code":
      return {
        ...state,
        code: value,
        name: state.name === state.code ? value : state.name
      };
    case "name":
      return {
        ...state,
        name: value
      };
    case "timeLimit":
      return {
        ...state,
        timeLimit: parseInt(value) || 0
      };
    case "memoryLimit":
      return {
        ...state,
        memoryLimit: parseInt(value) || 0
      };
    default:
      return state;
  }
}

function useProblemForm(initialProblem: Problem = defaultProblem) {
  let [problem, dispatch] = useReducer(reducer, initialProblem);

  let inputProps = useCallback(
    function(name: ActionName) {
      return {
        value: problem[name],
        onChange: function(e: ChangeEvent) {
          let { value } = e.target as HTMLInputElement;
          dispatch({ name, value });
        },
        required: true
      };
    },
    [problem]
  );

  let setDefaultCodeName = useCallback(
    function(value: string) {
      if (problem.code === "") {
        dispatch({ name: "code", value });
      }
      if (problem.name === "") {
        dispatch({ name: "name", value });
      }
    },
    [problem.code, problem.name]
  );

  return {
    problem,
    setDefaultCodeName,
    inputProps
  };
}

export default useProblemForm;
