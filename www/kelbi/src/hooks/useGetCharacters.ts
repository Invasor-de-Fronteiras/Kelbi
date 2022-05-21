import { useEffect, useState } from 'react';
import {
  Character,
  getCharacters,
  getLastAuthResult,
  getSignResult,
  LastAuthResult,
  SignResult,
} from '../utils/launcher';

interface GetCharacterHook {
  loading: boolean;
  characters: Character[];
  isNewAccount: boolean;
}

export function useGetCharacters(): GetCharacterHook {
  const [loading, setLoading] = useState(true);
  const [characters, setCharacters] = useState<Character[]>([]);
  const [isNewAccount, setIsNewAccount] = useState(false);

  useEffect(() => {
    const interval = setInterval(() => {
      const lastAuth = getLastAuthResult();
      const signRes = getSignResult();
      const chars = getCharacters();

      if (lastAuth === LastAuthResult.InLoading) {
        setLoading(true);
        return;
      }

      clearInterval(interval);

      if (
        lastAuth === LastAuthResult.AuthSuccess &&
        signRes === SignResult.SignSuccess &&
        chars.length > 0
      ) {
        setLoading(false);
        setCharacters(chars);
      }
    }, 100);

    return () => clearInterval(interval);
  }, []);

  return { loading, characters, isNewAccount };
}
