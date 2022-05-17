import { useLauncher } from '../context/LauncherContext';
import { useMutation } from './useMutation';

export interface LoginInput {
  username: string;
  password: string;
  rememberMe: boolean;
}

export function useLogin() {
  const { setLoggedIn } = useLauncher();

  return useMutation(({ username, password, rememberMe }: LoginInput): void => {
    //@ts-ignore
    window.external.loginCog(username + '+', password, password);

    localStorage.setItem('rememberMe', String(rememberMe));

    if (rememberMe) {
      localStorage.setItem('username', username);
      localStorage.setItem('password', password);
    }

    setLoggedIn(true);
  });
}
