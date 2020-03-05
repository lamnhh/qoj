import React, { useEffect, useState } from "react";
import Contest from "../models/Contest";
import moment from "moment";
import { useHistory } from "react-router-dom";

function ContestCountDown({ contest }: { contest: Contest }) {
  let [time, setTime] = useState(moment(contest.startDate).diff(moment()));
  let history = useHistory();

  useEffect(
    function() {
      let start = moment(contest.startDate);
      let timer = setInterval(function() {
        if (start.isBefore(moment())) {
          clearInterval(timer);
          history.push("/contest/" + contest.id);
          return;
        }
        setTime(start.diff(moment()));
      }, 1000);

      return function() {
        clearInterval(timer);
      };
    },
    [contest.startDate, contest.id, history]
  );

  return (
    <div className="count-down">
      <h2 className="count-down__header">Before the contest</h2>
      <h2 className="count-down__content">
        {moment.duration(time).format("DD[d] HH[h] mm[m] ss[s]")}
      </h2>
    </div>
  );
}

export default ContestCountDown;
