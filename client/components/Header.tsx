import React, { useContext } from "react";
import { NavLink, Link } from "react-router-dom";
import AppContext from "../contexts/AppContext";

function Header() {
  let { user } = useContext(AppContext);
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
        </ul>
      </nav>
      <nav>
        {user !== null ? (
          <ul>
            <li>{user.username}</li>
            <li>
              <Link to="/logout">Sign out</Link>
            </li>
          </ul>
        ) : (
          <ul>
            <li>
              <Link to="/login">Sign in</Link>
            </li>
            <li>
              <Link to="/register">Register</Link>
            </li>
          </ul>
        )}
      </nav>
    </header>
  );
}

export default Header;
