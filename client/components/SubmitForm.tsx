import React, {
  useCallback,
  FormEvent,
  ChangeEvent,
  useRef,
  useState,
  useEffect
} from "react";
import { useHistory } from "react-router-dom";
import request from "../helpers/request";
import Language from "../models/Language";
import CodeMirror, { EditorFromTextArea } from "codemirror";
import Problem from "../models/Problem";
import "codemirror/mode/clike/clike";

interface FormElements extends HTMLFormElement {
  languageId: HTMLSelectElement;
  problemId: HTMLSelectElement | HTMLInputElement;
  file: HTMLInputElement;
}

function SubmitForm({
  problemList,
  redirectUrl
}: {
  problemList: Array<Problem>;
  redirectUrl: string;
}) {
  let history = useHistory();
  let codeRef = useRef<HTMLTextAreaElement | null>(null);
  let editor = useRef<EditorFromTextArea | null>(null);

  useEffect(function() {
    editor.current = CodeMirror.fromTextArea(codeRef.current!, {
      lineNumbers: true,
      mode: "text/x-c++src"
    });
  }, []);

  let [languageList, setLanguagelist] = useState<Array<Language>>([]);
  useEffect(function() {
    request("/api/language").then(setLanguagelist);
  }, []);

  let handleSubmit = useCallback(
    function(event: FormEvent<HTMLFormElement>) {
      event.preventDefault();
      let form = event.target as FormElements;
      let code = editor.current!.getValue();
      let languageId = parseInt(form.languageId.value);
      let problemId = parseInt(form.problemId.value);

      request("/api/submission", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          problemId,
          languageId,
          code
        })
      })
        .then(function() {
          history.push(redirectUrl);
        })
        .catch(function({ error }) {
          alert(error);
        });
    },
    [redirectUrl, history]
  );

  let onFileUpload = useCallback(function(e: ChangeEvent) {
    let element = e.target as HTMLInputElement;
    let files = element.files;
    if (files === null) {
      return;
    }

    let file = files[0];
    element.value = "";
    if (file.size > 50000) {
      alert("Solution file exceeds 50000B");
      return;
    }

    let reader = new FileReader();
    reader.onload = function() {
      editor.current!.setValue(String(this.result));
    };

    reader.readAsText(file, "utf-8");
  }, []);

  return (
    <form className="submit-form" onSubmit={handleSubmit}>
      <label>
        <span className="submit-form__field-name">Language</span>
        <select className="submit-form__input" name="languageId">
          {languageList.map(function(language) {
            return (
              <option key={language.id} value={language.id}>
                {language.name}
              </option>
            );
          })}
        </select>
      </label>
      {problemList.length === 1 ? (
        <input type="hidden" name="problemId" value={problemList[0].id}></input>
      ) : (
        <label>
          <span className="submit-form__field-name">Problem</span>
          <select
            className="submit-form__input"
            name="problemId"
            placeholder="Choose problem">
            {problemList.map(function(problem, idx) {
              return (
                <option key={problem.id} value={problem.id}>
                  {String.fromCharCode("A".charCodeAt(0) + idx)}. {problem.name}
                </option>
              );
            })}
          </select>
        </label>
      )}
      <label>
        <span className="submit-form__field-name">Source code</span>
        <textarea className="submit-form__input" ref={codeRef}></textarea>
      </label>
      <label>
        <span className="submit-form__field-name">Or choose file</span>
        <input
          className="submit-form__input"
          type="file"
          name="file"
          onChange={onFileUpload}
        />
      </label>
      <button type="submit" className="submit-form__submit-btn">
        Submit
      </button>
    </form>
  );
}

export default SubmitForm;
