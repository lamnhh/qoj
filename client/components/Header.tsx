import React, {
  useContext,
  useState,
  useCallback,
  useEffect,
  useRef
} from "react";
import { Link, matchPath, useLocation } from "react-router-dom";
import AppContext from "../contexts/AppContext";

function MainNavItem({
  path,
  label,
  onClick
}: {
  path: string;
  label: string;
  onClick: () => void;
}) {
  let pathname = useLocation().pathname;
  let [isFocussed, setIsFocussed] = useState(false);
  let match = matchPath(pathname, { path, exact: true });

  return (
    <li className={(isFocussed ? "active " : "") + (match ? "match" : "")}>
      <Link
        to={path}
        onClick={onClick}
        onFocus={() => setIsFocussed(true)}
        onBlur={() => setIsFocussed(false)}>
        {label}
      </Link>
    </li>
  );
}

function Header() {
  let { user } = useContext(AppContext);

  let [active, setActive] = useState(false);
  let disableNavbar = useCallback(function() {
    setActive(false);
  }, []);

  let navbarRef = useRef<HTMLElement | null>(null);
  useEffect(
    function() {
      function handleClick(e: Event) {
        let el = e.target as Node;
        if (active && !navbarRef.current!!.contains(el)) {
          setActive(false);
        }
      }
      window.addEventListener("click", handleClick);
      return function() {
        window.removeEventListener("click", handleClick);
      };
    },
    [active]
  );

  return (
    <header className="header">
      <div className="align-left-right header__main-nav--wrapper">
        <Link to="/" className="header__logo-wrapper">
          <h1 className="header__logo">QHH Online Judge</h1>
        </Link>
        <button
          type="button"
          id="navbar-toggler"
          onClick={() => setActive(p => !p)}>
          <img
            src="/static/images/navbar-toggler.svg"
            alt="Toggle navigation bar"></img>
        </button>
        <nav
          className={"header__main-nav " + (active ? "active" : "inactive")}
          ref={navbarRef}>
          <ul>
            <MainNavItem path="/" label="Problemset" onClick={disableNavbar} />
            <MainNavItem
              path="/status"
              label="Submission"
              onClick={disableNavbar}
            />
            <MainNavItem
              path="/contest"
              label="Contests"
              onClick={disableNavbar}
            />
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
