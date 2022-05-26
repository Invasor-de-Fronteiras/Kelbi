/* eslint-disable @typescript-eslint/no-explicit-any */
import { useState } from 'react';

interface UsePromiseState<T = any> {
  isLoading: boolean;
  error: Error | null;
  data: T | null;
}

const defaultState: UsePromiseState = {
  data: null,
  error: null,
  isLoading: false,
};

export function usePromise<Params = any, Result = any>(
  handler: (...args: Params[]) => Promise<Result>,
  onResolved?: (data: Result) => void,
  onRejected?: (error: Error) => void,
) {
  const [state, setState] = useState<UsePromiseState<Result>>(defaultState);

  function mutate(...args: Params[]) {
    setState((state) => ({ error: null, data: state.data, isLoading: true }));

    handler(...args)
      .then((data) => {
        onResolved?.(data);
        setState({
          error: null,
          isLoading: false,
          data,
        });
      })
      .catch((error) => {
        onRejected?.(error);
        setState({
          error,
          isLoading: false,
          data: null,
        });
      });
  }

  return { mutate, ...state };
}
