import { createNewCharacter, startGame } from '../utils/launcher';
import { usePromise } from './common/usePromise';

export function useCreateCharacter() {
  return usePromise(async () => {
    const newChar = await createNewCharacter();
    startGame(newChar.uid);
  });
}
