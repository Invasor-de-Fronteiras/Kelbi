import TerrariaSong from '../assets/audio/terraria.mp3';
import SteveSong from '../assets/audio/steve.mp3';
import UndertaleSong from '../assets/audio/undertale.mp3';
import RobloxDeathSong from '../assets/audio/roblox_death.mp3';

import { probably, randomArr } from './util';

export function playClickSong() {
  const audio = new Audio(SteveSong);
  audio.play();
}

export function playLoginSong() {
  const audio = new Audio(TerrariaSong);
  audio.play();
}

export function randomSong() {
  if (probably(5)) {
    const song = randomArr([RobloxDeathSong, UndertaleSong]);
    const audio = new Audio(song);
    audio.play();
  }
}

export function playDigitSong() {
  randomSong();
}
