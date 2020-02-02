import React from "react";
import { range } from "../helpers/common-helper";

interface PaginationButtonProps {
  page: number | string;
  active: boolean;
  disabled?: boolean;
  onClick: () => void;
}

interface PaginationProps {
  totalCount: number;
  pageSize: number;
  currentPage: number;
  onPageChange: (page: number) => void;
}

function PaginationButton({
  page,
  active,
  onClick,
  disabled
}: PaginationButtonProps) {
  return (
    <button
      type="button"
      className={`${active ? "active " : ""}pagination-button`}
      disabled={disabled}
      onClick={onClick}>
      {page}
    </button>
  );
}

PaginationButton.defaultValue = {
  disabled: false
};

function Pagination({
  totalCount,
  pageSize,
  currentPage,
  onPageChange
}: PaginationProps) {
  let pageCount = Math.ceil(totalCount / pageSize);
  return (
    <div className="pagination">
      <PaginationButton
        page="«"
        active={false}
        disabled={currentPage === 1}
        onClick={function() {
          onPageChange(Math.max(1, currentPage - 1));
        }}
      />
      {range(1, pageCount).map(function(idx) {
        return (
          <PaginationButton
            page={idx}
            active={currentPage === idx}
            onClick={function() {
              onPageChange(idx);
            }}
          />
        );
      })}
      <PaginationButton
        page="»"
        active={false}
        disabled={currentPage === pageCount}
        onClick={function() {
          onPageChange(Math.min(pageCount, currentPage + 1));
        }}
      />
    </div>
  );
}

export default Pagination;
