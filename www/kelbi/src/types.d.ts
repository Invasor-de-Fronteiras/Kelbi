declare module '*.png';
declare module '*.jpg';
declare module '*.gif';
declare module '*.svg';

declare module global {
  interface Window {
    external: LauncherFunctions;
  }
}

interface LauncherFunctions {
  playSound(song: LauncherSongs): void;
  beginDrag(active: boolean): void;
  openBrowser(url: string): void;
  restartMhf(): void;
  minimizeWindow(): void;
  closeWindow(): void;
}

type LauncherSongs = 'IDR_WAV_SEL' | 'IDR_WAV_OK' | 'IDR_WAV_PRE_LOGIN' | 'IDR_NIKU';
