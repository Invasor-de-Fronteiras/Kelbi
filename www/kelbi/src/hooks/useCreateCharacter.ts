import { useState } from 'react';
import { createNewCharacter, startGame } from '../utils/launcher';

export function useCreateCharacter() {
  const [state, setState] = useState({
    isLoading: false,
    error: null,
  });

  const mutate = () => {
    setState({ isLoading: true, error: null });
    createNewCharacter(
      (char) => startGame(char.uid),
      (error) => {
        setState({ isLoading: false, error });
      },
    );
  };

  return {
    mutate,
    ...state,
  };
}
