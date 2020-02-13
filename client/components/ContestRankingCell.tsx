import React from "react";

function ContestRankingCell({ score }: { score?: number; maxScore: number }) {
  return <td className="score">{score}</td>;
}

export default ContestRankingCell;
