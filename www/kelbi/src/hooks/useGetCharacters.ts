import type { Character } from '../utils/launcher';

import { useEffect, useState } from 'react';

import { getCharacters, isNeAccountChar, LastAuthResult, SignResult } from '../utils/launcher';
import { removeItem } from '../utils/util';

interface GetCharacterHook {
  loading: boolean;
  characters: Character[];
  isNewAccount: boolean;
  newAccountUID: string;
}

export function useGetCharacters(): GetCharacterHook {
  const [state, setState] = useState<Omit<GetCharacterHook, 'isNewAccount'>>({
    loading: true,
    characters: [],
    newAccountUID: '',
  });

  useEffect(() => {
    const interval = setInterval(() => {
      const lastAuth = window.external.getLastAuthResult();
      const signRes = window.external.getSignResult();
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

        const characters = removeItem(chars, (char) => {
          if (isNeAccountChar(char)) {
            newAccountUID = char.uid;
            return true;
          }

          return false;
        });

        setState({
          loading: false,
          characters,
          newAccountUID,
        });
      }
    }, 100);

    return () => clearInterval(interval);
  }, []);

  return {
    ...state,
    isNewAccount: state.characters.length === 0,
  };
}
