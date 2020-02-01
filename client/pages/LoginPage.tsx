import React from "react";
import LoginForm from "../components/LoginForm";

function LoginPage() {
  return (
    <>
      <header className="page-name align-left-right">
        <h1>Sign In</h1>
      </header>
      <section className="auth-page">
        <LoginForm />
      </section>
    </>
  );
}

export default LoginPage;
