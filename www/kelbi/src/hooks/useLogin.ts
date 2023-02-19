import { useState } from 'react';
import { useLauncher } from '../context/LauncherContext';
import { LocaleError } from '../i18n/LocaleError';
import { useTranslate } from '../i18n/useTranslate';
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
  const { t } = useTranslate();

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

        if (lastAuth === LastAuthResult.None) {
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
          handleError(new LocaleError('login_incorrect_password'));
        } else if (LastAuthResult.AuthErrorNet) {
          handleError(new LocaleError('server_offline_error'));
        } else {
          handleError(new LocaleError('login_unknown_error'));
        }
      }, 100);
    } catch (err) {
      handleError(err);
    }
  };

  return { ...state, mutate, cleanError: () => setState({ ...state, error: null }) };
}
