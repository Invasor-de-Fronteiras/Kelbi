import { links } from '../constants';

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
