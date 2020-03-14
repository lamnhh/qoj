import { Dispatch } from "react";

type DispatchWrapper<T> = Dispatch<T | ((prevState: T) => T)>;

export default DispatchWrapper;
