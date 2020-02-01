import React, { useCallback, FormEvent } from "react";
import { useHistory } from "react-router-dom";
import request from "../helpers/request";

interface FormElements extends HTMLFormElement {
  file: HTMLInputElement;
}

function SubmitForm({ problemId }: { problemId: string }) {
  let history = useHistory();
  let handleSubmit = useCallback(
    function(event: FormEvent<HTMLFormElement>) {
      event.preventDefault();
      let form = event.target as FormElements;
      let file = form.file.files![0];

      let body = new FormData();
      body.append("problemId", problemId);
      body.append("file", file);

      request("/api/submission", {
        method: "POST",
        body
      }).then(function() {
        history.push("/status");
      });
    },
    [problemId]
  );

  return (
    <form className="submit-form" onSubmit={handleSubmit}>
      <label>
        <span>Language</span>
        <select className="submit-form__language">
          <option key="cpp">C</option>
          <option key="cpp">C++</option>
          <option key="pas">Pascal</option>
          <option key="jav">Java</option>
        </select>
      </label>
      <label>
        <span>Source code</span>
        <textarea className="submit-form__editor"></textarea>
      </label>
      <label>
        <span>Or choose file</span>
        <input type="file" name="file" required />
      </label>
      <button type="submit">Submit</button>
    </form>
  );
}

export default SubmitForm;
