import React, { PropsWithChildren, useContext } from "react";
import AdminHeader from "./AdminHeader";
import AdminNavigation from "./AdminNavigation";
import AppContext from "../contexts/AppContext";
import { Redirect } from "react-router-dom";

interface AdminLayoutProps {
  title: string;
  requireLogin?: boolean;
}

function AdminLayout({
  title,
  requireLogin = true,
  children
}: PropsWithChildren<AdminLayoutProps>) {
  let { user } = useContext(AppContext);
  if (requireLogin && !user) {
    return <Redirect to="/" />;
  }
  return (
    <div className="align-left-right">
      <AdminHeader />
      <AdminNavigation title={title} loggedIn={!!user} />
      {children}
    </div>
  );
}

export default AdminLayout;
