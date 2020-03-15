import React, { useCallback, FormEvent, useState, useEffect } from "react";
import useSWR from "swr";
import DateTime from "react-datetime";
import request from "../helpers/request";
import Problem from "../models/Problem";
import useSelectedProblemList from "../hooks/useSelectedProblemList";
import moment, { Moment } from "moment";
import Contest from "../models/Contest";
import AdminContest from "../models/AdminContest";

interface AdminContestFormProps {
  action: string;
  defaultContest: Contest;
  defaultSelected: number[];
  handleSubmit: (contest: AdminContest) => void;
}

interface AddFormElement extends HTMLFormElement {
  problemId: HTMLSelectElement;
}

interface ContestFormElement extends HTMLFormElement {
  contestName: HTMLInputElement;
  duration: HTMLInputElement;
}

function AdminContestForm({
  action,
  defaultContest,
  defaultSelected,
  handleSubmit
}: AdminContestFormProps) {
  let { data } = useSWR("/api/problem", request);

  // prettier-ignore
  let { selectedValues, isSelected, select, unselect } = useSelectedProblemList<number>(defaultSelected);

  let onAddProblem = useCallback(
    function(e: FormEvent) {
      e.preventDefault();
      let form = e.target as AddFormElement;
      let problemId = parseInt(form.problemId.value);
      if (!isNaN(problemId)) {
        select(problemId);
        form.problemId.value = "";
      }
    },
    [select]
  );

  let [date, setDate] = useState<Moment>(() => moment());
  useEffect(
    function() {
      if (defaultContest.id === 0) {
        setDate(moment());
      } else {
        setDate(moment(defaultContest.startDate));
      }
    },
    [defaultContest]
  );

  let onSubmit = useCallback(
    function(e: FormEvent) {
      e.preventDefault();
      let form = e.target as ContestFormElement;

      if (selectedValues.length === 0) {
        alert("Please choose at least 1 problem");
        return;
      }

      let name = form.contestName.value;
      let startDate = date;
      let duration = parseInt(form.duration.value) ?? 120;
      let problemList = selectedValues;

      handleSubmit({
        name,
        startDate: startDate.toDate(),
        problemList,
        duration
      });
    },
    [date, selectedValues, handleSubmit]
  );

  if (!data) {
    return <section className="loading-msg">Loading</section>;
  }

  let problemList: Problem[] = data;
  return (
    <section className="contest-form--wrapper">
      <div className="contest-form">
        <form
          className="contest-form__metadata"
          id="metadata-form"
          onSubmit={onSubmit}>
          <label>
            <span className="contest-form__field-name">Contest name</span>
            <div className="contest-form__field-value">
              <input
                type="text"
                name="contestName"
                defaultValue={defaultContest.name}
                required
                autoFocus
              />
            </div>
          </label>
          <label>
            <span className="contest-form__field-name">Start time</span>
            <div className="contest-form__field-value">
              <DateTime
                value={date}
                onChange={function(date) {
                  setDate(moment(date));
                }}></DateTime>
            </div>
          </label>
          <label>
            <span className="contest-form__field-name">Duration</span>
            <div className="contest-form__field-value">
              <input
                type="text"
                name="duration"
                defaultValue={defaultContest.duration}
                inputMode="numeric"
                pattern="[0-9]*"
                required
              />
              <p className="contest-form__field-desc">Duration in minutes</p>
            </div>
          </label>
        </form>
        <form onSubmit={onAddProblem}>
          <table className="admin-table contest-form__selected">
            <thead>
              <tr>
                <th>#</th>
                <th>Problem Name</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {selectedValues.map(function(problemId) {
                let problem = problemList.find(({ id }) => id === problemId)!;
                return (
                  <tr key={problem.id}>
                    <td>{problem.id}</td>
                    <td>
                      {problem.code === problem.name
                        ? problem.code
                        : problem.code + " - " + problem.name}
                    </td>
                    <td className="contest-form__selected__action">
                      <button
                        type="button"
                        onClick={function() {
                          unselect(problem.id);
                        }}>
                        Discard
                      </button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
            <tfoot>
              <tr>
                <td></td>
                <td>
                  <select name="problemId" required>
                    <option value="" disabled selected>
                      --- Choose a problem ---
                    </option>
                    {problemList.map(function(problem) {
                      if (isSelected(problem.id)) {
                        return null;
                      }
                      return (
                        <option key={problem.id} value={problem.id}>
                          {problem.code === problem.name
                            ? problem.code
                            : problem.code + " - " + problem.name}
                        </option>
                      );
                    })}
                  </select>
                </td>
                <td className="contest-form__selected__action">
                  <button type="submit">Add</button>
                </td>
              </tr>
            </tfoot>
          </table>
        </form>
      </div>
      <div className="contest-form__submit-btn">
        <button type="submit" form="metadata-form">
          {action}
        </button>
      </div>
    </section>
  );
}

export default AdminContestForm;
