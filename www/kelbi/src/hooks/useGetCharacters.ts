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
  newAccountUID: string;
}

export function useGetCharacters(): GetCharacterHook {
  const [state, setState] = useState<GetCharacterHook>({
    loading: true,
    characters: [],
    isNewAccount: false,
    username: '',
    newAccountUID: '',
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
        let newAccountUID = '';

        let characters = removeItem(chars, (char) => {
          if (isNeAccountChar(char)) {
            newAccountUID = char.uid;
            return true;
          }

          return false;
        });

        setState({
          loading: false,
          characters,
          username: getUserId(),
          isNewAccount: characters.length === 0,
          newAccountUID,
        });
      }
    }, 100);

    return () => clearInterval(interval);
  }, []);

  return state;
}
