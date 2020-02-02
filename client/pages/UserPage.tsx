import React, { useState, useEffect } from "react";
import User, { emptyUser } from "../models/User";
import { useParams } from "react-router-dom";
import request from "../helpers/request";
import UserPageProblemList from "../components/UserPageProblemList";

interface UserPageRouterProps {
  username: string;
}

function UserPage() {
  let username = useParams<UserPageRouterProps>().username;
  let [user, setUser] = useState<User>(emptyUser);
  useEffect(
    function() {
      request(`/api/user/${username}/public`)
        .then(setUser)
        .catch(console.log);
    },
    [username]
  );

  return (
    <>
      <header className="page-name align-left-right">
        <h1>Profile of {username}</h1>
      </header>
      <section className="user-page align-left-right">
        <div className="user-page__info">
          <img
            src="/static/profile-picture-placeholder.png"
            alt={`${username}'s profile picture`}
          />
          <h1>{user.fullname}</h1>
          <h2>{user.username}</h2>
        </div>
        <div className="user-page__prob-list--wrapper">
          <UserPageProblemList
            url={`/api/user/${username}/solved`}
            title="Solved problems"
          />
          <UserPageProblemList
            url={`/api/user/${username}/partial`}
            title="Partially solved problems"
          />
        </div>
      </section>
    </>
  );
}

export default UserPage;