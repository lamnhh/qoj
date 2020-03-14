import React from "react";
import { NavLink, Link } from "react-router-dom";

interface AdminNavigationProps {
  title: string;
  loggedIn: boolean;
}

function AdminNavigation({ title, loggedIn }: AdminNavigationProps) {
  return (
    <nav className="navigation">
      <h1 className="page-title">{title}</h1>
      {loggedIn && (
        <ul>
          <li>
            <NavLink to="/search">Search</NavLink>
          </li>
          <span>|</span>
          <li>
            <NavLink to="/problem/new">New Problem</NavLink>
          </li>
          <li>
            <NavLink to="/problem" exact>
              View Problems
            </NavLink>
          </li>
          <span>|</span>
          <li>
            <NavLink to="/contest/new">New Contest</NavLink>
          </li>
          <li>
            <NavLink to="/contest" exact>
              View Contests
            </NavLink>
          </li>
          <span>|</span>
          <li>
            <Link to="/logout">Logout</Link>
          </li>
        </ul>
      )}
    </nav>
  );
}

export default AdminNavigation;
