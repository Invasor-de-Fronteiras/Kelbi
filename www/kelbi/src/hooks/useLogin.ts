import { useEffect, useState } from 'react';
import { useLauncher } from '../context/LauncherContext';
import { getLastAuthResult, getSignResult, LastAuthResult, SignResult } from '../utils/launcher';

export interface LoginInput {
  username: string;
  password: string;
  rememberMe: boolean;
}

export function useLogin() {
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

    try {
      // TESTAR SEM + OU COM TERCEIRO INPUT INCORRETO
      // TESTEI, MAS SEM RESULTADOS
      //@ts-ignore
      window.external.loginCog(input.username, input.password, input.password);
    } catch (err) {
      //@ts-ignore
      setState({ isLoading: false, error: err });
    }
  };

  useEffect(() => {
    if (!state.isLoading) return;

    const interval = setInterval(() => {
      const lastAuth = getLastAuthResult();
      const signRes = getSignResult();

      if (lastAuth === LastAuthResult.None || signRes === SignResult.None) {
        return;
      }

      if (lastAuth === LastAuthResult.InLoading) {
        setState({ isLoading: true, error: null });
        setIsLoading(true);
      } else if (lastAuth === LastAuthResult.AuthSuccess && signRes === SignResult.SignSuccess) {
        setLoggedIn(true);
        setIsLoading(false);
        setState({ isLoading: false, error: null });
      } else if (signRes === SignResult.NotMatchPassword) {
        setState({ isLoading: false, error: new Error('senha incorreta!') });
        setIsLoading(false);
      } else {
        setState({
          isLoading: false,
          error: new Error(`falha na autenticação! ${lastAuth} ${signRes}`),
        });
        setIsLoading(false);
      }
    }, 100);

    return () => clearInterval(interval);
  }, [state]);

  return { ...state, mutate, cleanError: () => setState({ ...state, error: null }) };
}
