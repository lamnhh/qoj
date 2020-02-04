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

interface FormElements extends HTMLFormElement {
  languageId: HTMLSelectElement;
  file: HTMLInputElement;
}

function SubmitForm({ problemId }: { problemId: string }) {
  let history = useHistory();
  let codeRef = useRef<HTMLTextAreaElement>(null);
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
