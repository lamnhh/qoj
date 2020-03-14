import React, { FormEvent } from "react";
import { useHistory } from "react-router-dom";

interface SearchFormElement extends HTMLFormElement {
  search: HTMLInputElement;
}

function AdminSearchPage() {
  let history = useHistory();

  function onSearch(e: FormEvent) {
    e.preventDefault();
    let form = e.target as SearchFormElement;

    let search = form.search.value;
    history.push("/problem?search=" + search);
  }

  return (
    <form className="search-form" onSubmit={onSearch}>
      <input type="text" name="search" autoFocus></input>
      <p className="search-form__desc">
        Try examples: &quot;NK&quot;, &quot;V11&quot;, etc
      </p>
      <button type="submit">Search</button>
    </form>
  );
}

export default AdminSearchPage;
