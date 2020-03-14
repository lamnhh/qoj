import { useState, useCallback, useEffect } from "react";

function useSelectedProblemList<T>(defaultValues: Array<T> = []) {
  let [selected, setSelected] = useState(function() {
    return new Set<T>(defaultValues);
  });

  useEffect(
    function() {
      setSelected(new Set<T>(defaultValues));
    },
    [defaultValues]
  );

  let isSelected = useCallback(
    function(id: T) {
      return selected.has(id);
    },
    [selected]
  );

  let select = useCallback(function(id: T) {
    setSelected(function(prev) {
      let next = new Set(prev);
      next.add(id);
      return next;
    });
  }, []);

  let unselect = useCallback(function(id: T) {
    setSelected(function(prev) {
      let next = new Set(prev);
      next.delete(id);
      return next;
    });
  }, []);

  return {
    selectedValues: Array.from(selected),
    isSelected,
    select,
    unselect
  };
}

export default useSelectedProblemList;
