import React, { useEffect, useState, useContext } from "react";
import { clearToken } from "../helpers/auth";
import AppContext from "../contexts/AppContext";
import { Redirect } from "react-router-dom";

interface LogoutProps {
  logoutUrl: string;
  returnUrl: string;
}

function Logout({ logoutUrl, returnUrl }: LogoutProps) {
  let [running, setRunning] = useState(true);
  let { setUser } = useContext(AppContext);

  useEffect(
    function() {
      fetch(logoutUrl)
        .then(function(res) {
          if (res.ok) {
            clearToken();
            setUser(null);
          }
        })
        .finally(function() {
          setRunning(false);
        });
    },
    [logoutUrl, setUser]
  );

  if (running) {
    return null;
  }

  return <Redirect to={returnUrl} />;
}

export default Logout;
