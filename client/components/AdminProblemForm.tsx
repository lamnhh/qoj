import React, { useRef, useCallback, FormEvent, ChangeEvent } from "react";
import useProblemForm from "../hooks/useProblemForm";
import Problem from "../models/Problem";
import AdminProblem from "../models/AdminProblem";

interface AdminProblemFormProps {
  title: string;
  action: string;
  defaultProblem: Problem;
  requireFile: boolean;
  handleSubmit: (problem: AdminProblem) => void;
}

function AdminProblemForm({
  title,
  action,
  defaultProblem,
  requireFile,
  handleSubmit
}: AdminProblemFormProps) {
  let { problem, inputProps, setDefaultCodeName } = useProblemForm(
    defaultProblem
  );
  let testRef = useRef<HTMLInputElement>(null);

  let onSubmit = useCallback(
    function(e: FormEvent) {
      e.preventDefault();

      handleSubmit({
        code: problem.code,
        name: problem.name,
        timeLimit: problem.timeLimit / 1000,
        memoryLimit: problem.memoryLimit,
        file: testRef.current?.files?.[0]
      });
    },
    [problem, handleSubmit]
  );

  let onUpload = useCallback(
    function(e: ChangeEvent) {
      let input = e.target as HTMLInputElement;
      if (input.files && input.files.length > 0) {
        let file = input.files[0];
        let name = file.name
          .split(".")
          .slice(0, -1)
          .join(".");
        setDefaultCodeName(name);
      }
    },
    [setDefaultCodeName]
  );

  return (
    <section className="problem-form--wrapper">
      <div className="problem-form">
        <h2 className="problem-form__title">{title}</h2>
        <form onSubmit={onSubmit}>
          <label>
            <span className="problem-form__field-name">Code</span>
            <div className="problem-form__field-value">
              <input type="text" {...inputProps("code")} autoFocus></input>
              <p className="problem-form__field-desc">
                Problem identifier (i.e. NKPALIN, V11STR)
              </p>
            </div>
          </label>
          <label>
            <span className="problem-form__field-name">Name</span>
            <div className="problem-form__field-value">
              <input type="text" {...inputProps("name")}></input>
              <p className="problem-form__field-desc">
                Problem name (i.e. Chuỗi đối xứng, Tìm xâu)
              </p>
            </div>
          </label>
          <label>
            <span className="problem-form__field-name">Time limit (ms)</span>
            <div className="problem-form__field-value">
              <input type="text" {...inputProps("timeLimit")}></input>
              <p className="problem-form__field-desc not-show">placeholder</p>
            </div>
          </label>
          <label>
            <span className="problem-form__field-name">Memory limit (MB)</span>
            <div className="problem-form__field-value">
              <input type="text" {...inputProps("memoryLimit")}></input>
              <p className="problem-form__field-desc not-show">placeholder</p>
            </div>
          </label>
          <label>
            <span className="problem-form__field-name">Testdata</span>
            <div className="problem-form__field-value">
              <input
                type="file"
                name="test"
                required={requireFile}
                ref={testRef}
                onChange={onUpload}
                accept=".zip"></input>
              <p className="problem-form__field-desc">
                Zipped test data (Themis format)
              </p>
            </div>
          </label>
          <button type="submit">{action}</button>
        </form>
      </div>
    </section>
  );
}

export default AdminProblemForm;
