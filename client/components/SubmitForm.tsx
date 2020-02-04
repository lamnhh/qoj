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

interface FormElements extends HTMLFormElement {
  languageId: HTMLSelectElement;
  code: HTMLTextAreaElement;
  file: HTMLInputElement;
}

function SubmitForm({ problemId }: { problemId: string }) {
  let history = useHistory();
  let codeRef = useRef<HTMLTextAreaElement>(null);

  let [languageList, setLanguagelist] = useState<Array<Language>>([]);
  useEffect(function() {
    request("/api/language").then(setLanguagelist);
  }, []);

  let handleSubmit = useCallback(
    function(event: FormEvent<HTMLFormElement>) {
      event.preventDefault();
      let form = event.target as FormElements;
      let code = form.code.value;
      let languageId = parseInt(form.languageId.value);

      request("/api/submission", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          problemId: parseInt(problemId),
          languageId,
          code
        })
      }).then(function() {
        history.push("/status");
      });
    },
    [problemId]
  );

  let onFileUpload = useCallback(function(e: ChangeEvent) {
    let element = e.target as HTMLInputElement;
    let files = element.files;
    if (files === null) {
      return;
    }

    let file = files[0];
    if (file.size > 50000) {
      alert("Solution file exceeds 50000B");
      element.value = "";
      return;
    }

    let reader = new FileReader();
    reader.onload = function() {
      if (codeRef.current) {
        codeRef.current.value = String(this.result);
      }
    };

    reader.readAsText(file, "utf-8");
  }, []);

  return (
    <form className="submit-form" onSubmit={handleSubmit}>
      <label>
        <span>Language</span>
        <select className="submit-form__language" name="languageId">
          {languageList.map(function(language) {
            return (
              <option key={language.id} value={language.id}>
                {language.name}
              </option>
            );
          })}
        </select>
      </label>
      <label>
        <span>Source code</span>
        <textarea
          ref={codeRef}
          className="submit-form__editor"
          name="code"></textarea>
      </label>
      <label>
        <span>Or choose file</span>
        <input type="file" name="file" required onChange={onFileUpload} />
      </label>
      <button type="submit">Submit</button>
    </form>
  );
}

export default SubmitForm;
