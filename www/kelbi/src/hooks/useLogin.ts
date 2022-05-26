import { useState } from 'react';
import { useLauncher } from '../context/LauncherContext';
import { LastAuthResult, SignResult } from '../utils/launcher';

export interface LoginInput {
  username: string;
  password: string;
  autoLogin: boolean;
}

interface LoginHookProps {
  onSuccess?: (input: LoginInput) => void;
}

export function useLogin({ onSuccess = () => null }: LoginHookProps) {
  const { setIsLoading, setLoggedIn } = useLauncher();
  const [state, setState] = useState<{
    isLoading: boolean;
    error: Error | null;
  }>({
    isLoading: false,
    error: null,
  });

  const mutate = (input: LoginInput) => {
    setState({ isLoading: true, error: null });
    setIsLoading(true);

    const handleSuccess = () => {
      setLoggedIn(true);
      setIsLoading(false);
      setState({ isLoading: false, error: null });
      onSuccess(input);
    };

    const handleError = (error: Error) => {
      setState({ isLoading: false, error });
      setIsLoading(false);
    };

    try {
      window.external.loginCog(input.username, input.password, input.password);

      const interval = setInterval(() => {
        const lastAuth = window.external.getLastAuthResult();
        const signRes = window.external.getSignResult();

        if (lastAuth === LastAuthResult.None || signRes === SignResult.None) {
          return;
        }

        if (lastAuth === LastAuthResult.InLoading) {
          setState({ isLoading: true, error: null });
          setIsLoading(true);
          return;
        }

        clearInterval(interval);

        if (lastAuth === LastAuthResult.AuthSuccess && signRes === SignResult.SignSuccess) {
          handleSuccess();
        } else if (signRes === SignResult.NotMatchPassword) {
          handleError(new Error('senha incorreta!'));
        } else {
          handleError(new Error(`falha na autenticação! ${lastAuth} ${signRes}`));
        }
      }, 100);
    } catch (err) {
      handleError(err);
    }
  };

  return { ...state, mutate, cleanError: () => setState({ ...state, error: null }) };
}
