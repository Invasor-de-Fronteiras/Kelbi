import { useLauncher } from '../context/LauncherContext';
import { useMutation } from './useMutation';

export interface LoginInput {
  username: string;
  password: string;
  rememberMe: boolean;
}

export function useLogin() {
  return useMutation<Promise<void>, LoginInput>(
    async ({ username, password }: LoginInput): Promise<void> =>
      //@ts-ignore
      window.external.loginCog(username + '+', password, password),
    {
      onSuccess: (_, { password, rememberMe, username }) => {
        localStorage.setItem('rememberMe', String(rememberMe));

        if (rememberMe) {
          localStorage.setItem('username', username);
          localStorage.setItem('password', password);
        }

        // setLoggedIn(true);
      },
    },
  );
}
