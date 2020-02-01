import React from "react";
import RegisterForm from "../components/RegisterForm";

function RegisterPage() {
  return (
    <>
      <header className="page-name align-left-right">
        <h1>Register</h1>
      </header>
      <section className="auth-page">
        <RegisterForm />
      </section>
    </>
  );
}

export default RegisterPage;
