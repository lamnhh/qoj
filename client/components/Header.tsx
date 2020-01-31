import React from "react";
import { NavLink } from "react-router-dom";

function Header() {
  return (
    <header>
      <img src="/static/logo.jpeg" alt="Logo" />
      <nav>
        <ul>
          <li>
            <NavLink to="/">Problemset</NavLink>
          </li>
          <li>
            <NavLink to="/status">Submission</NavLink>
          </li>
          <li>
            <NavLink to="/login">Sign in</NavLink>
          </li>
          <li>
            <NavLink to="/register">Register</NavLink>
          </li>
        </ul>
      </nav>
    </header>
  );
}

export default Header;
