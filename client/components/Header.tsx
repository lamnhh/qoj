import React, { useContext, useState } from "react";
import { Link, matchPath, useLocation } from "react-router-dom";
import AppContext from "../contexts/AppContext";

function MainNavItem({ path, label }: { path: string; label: string }) {
  let pathname = useLocation().pathname;
  let [isFocussed, setIsFocussed] = useState(false);
  let match = matchPath(pathname, { path, exact: true });

  return (
    <li className={(isFocussed ? "active " : "") + (match ? "match" : "")}>
      <Link
        to={path}
        onFocus={() => setIsFocussed(true)}
        onBlur={() => setIsFocussed(false)}>
        {label}
      </Link>
    </li>
  );
}

function Header() {
  let { user } = useContext(AppContext);
  return (
    <header className="header">
      <div className="align-left-right header__main-nav--wrapper">
        <Link to="/" className="header__logo-wrapper">
          <h1 className="header__logo">QHH Online Judge</h1>
        </Link>
        <nav className="header__main-nav">
          <ul>
            <MainNavItem path="/" label="Problemset" />
            <MainNavItem path="/status" label="Submission" />
            <MainNavItem path="/contests" label="Contests" />
          </ul>
        </nav>
      </div>
      <nav className="align-left-right header__user-nav">
        {user !== null ? (
          <ul>
            <li>
              <Link to={"/user/" + user.username}>{user.username}</Link>
            </li>
            <li className="divider"></li>
            <li>
              <Link to="/settings/profile">Settings</Link>
            </li>
            <li className="divider"></li>
            <li>
              <Link to="/logout">Sign out</Link>
            </li>
          </ul>
        ) : (
          <ul>
            <li>
              <Link to="/login">Sign in</Link>
            </li>
            <li className="divider"></li>
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
