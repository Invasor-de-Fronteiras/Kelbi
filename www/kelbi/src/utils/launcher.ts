import { links } from '../constants';

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

function openBrowser(url: string) {
  // @ts-ignore
  window.external.openBrowser(url);
}

export function openDiscord() {
  openBrowser(links.discordUrl);
}

export function openGithub() {
  openBrowser(links.githubUrl);
}

export function openConfig() {
  // @ts-ignore
  window.external.openMhlConfig();
}

export function minimizeWindow() {
  // @ts-ignore
  window.external.minimizeWindow();
}

export function closeWindow() {
  // @ts-ignore
  window.external.closeWindow();
}
