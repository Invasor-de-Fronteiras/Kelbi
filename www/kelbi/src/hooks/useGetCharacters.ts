import { useEffect, useState } from 'react';
import {
  Character,
  getCharacters,
  getLastAuthResult,
  getSignResult,
  LastAuthResult,
  SignResult,
  isNeAccountChar,
  getUserId,
} from '../utils/launcher';
import { removeItem } from '../utils/util';

interface GetCharacterHook {
  loading: boolean;
  characters: Character[];
  isNewAccount: boolean;
  username: string;
}

export function useGetCharacters(): GetCharacterHook {
  const [state, setState] = useState<GetCharacterHook>({
    loading: true,
    characters: [],
    isNewAccount: false,
    username: '',
  });

  useEffect(() => {
    const interval = setInterval(() => {
      const lastAuth = getLastAuthResult();
      const signRes = getSignResult();
      const chars = getCharacters();

      if (lastAuth === LastAuthResult.InLoading) {
        return;
      }

      clearInterval(interval);

      if (
        lastAuth === LastAuthResult.AuthSuccess &&
        signRes === SignResult.SignSuccess &&
        chars.length > 0
      ) {
        let isNewAccount = false;
        let characters = removeItem(chars, (char) => {
          if (isNeAccountChar(char)) {
            isNewAccount = true;
            return true;
          }

          return false;
        });

        setState({
          loading: false,
          characters,
          isNewAccount,
          username: getUserId(),
        });
      }
    }, 100);

    return () => clearInterval(interval);
  }, []);

  return state;
}
