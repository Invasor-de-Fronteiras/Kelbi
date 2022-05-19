import { useCallback, useState } from 'react';
import { useLauncher } from '../context/LauncherContext';

interface MutationHook<T, V> {
  mutate: MutationFn<T, V>;
  isLoading: boolean;
  error: Error | null;
  data: T | null;
}

type MutationFn<T, V> = (variables?: V) => Promise<T>;
type MutationHookOptions<T, V> = {
  onError?: (error: Error) => void;
  onSuccess?: (data: T, variables: V) => void;
};

export function useMutation<T, V>(
  handler: MutationFn<T, V>,
  options: MutationHookOptions<T, V>,
): MutationHook<T, V> {
  const { setIsLoading, isLoading } = useLauncher();

  // TODO: implement useReducer
  const [state, setState] = useState({
    isLoading: false,
    error: null,
    data: null,
  });

  const mutate = useCallback<MutationFn<T, V>>(
    // @ts-ignore
    async (props: V) => {
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
        options?.onSuccess?.(data, props);

        return data;
      } catch (error) {
        // @ts-ignore
        options.onError?.(error);

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
    isLoading,
    mutate,
  };
}
