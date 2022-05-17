import { useCallback, useState } from 'react';
import { useLauncher } from '../context/LauncherContext';

interface MutationHook<T, V> {
  mutate: MutationFn<T, V>;
  isLoading: boolean;
  error: Error | null;
  data: T | null;
}

type MutationFn<T, V> = (variables?: V) => Promise<T>;

export function useMutation<T, V>(handler: MutationFn<T, V>): MutationHook<T, V> {
  const { setIsLoading } = useLauncher();

  // TODO: implement useReducer
  const [state, setState] = useState({
    isLoading: false,
    error: null,
    data: null,
  });

  const mutate = useCallback<MutationFn<T, V>>(
    // @ts-ignore
    async (props) => {
      setState({
        isLoading: true,
        error: null,
        data: null,
      });
      setIsLoading(true);

      try {
        const data = await handler(props);

        setState({
          isLoading: false,
          error: null,
          // @ts-ignore
          data,
        });

        return data;
      } catch (error) {
        setState({
          isLoading: false,
          //@ts-ignore
          error,
          data: null,
        });
      } finally {
        setIsLoading(false);
      }
    },
    [handler],
  );

  return {
    ...state,
    mutate,
  };
}
