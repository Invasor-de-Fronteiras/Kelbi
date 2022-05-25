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

export function openBrowser(url: string) {
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

export enum LastAuthResult {
  None = 'AUTH_NULL',
  AuthSuccess = 'AUTH_SUCCESS',
  InLoading = 'AUTH_PROGRESS',
  AuthErrorAcc = 'AUTH_ERROR_ACC',
}

export function getLastAuthResult(): LastAuthResult {
  // @ts-ignore
  return window.external.getLastAuthResult();
}

export enum SignResult {
  None = 'SIGN_UNKNOWN',
  SignSuccess = 'SIGN_SUCCESS',
  NotMatchPassword = 'SIGN_EPASS',
}

export function getSignResult(): SignResult {
  // @ts-ignore
  return window.external.getSignResult();
}

export function isEnableSessionId() {
  try {
    // @ts-ignore
    return window.external.isEnableSessionId();
  } catch {
    // IGNORANDO ERRO IGUAL LAUNCHER ORIGIONAL
    // TODO: DEBUGAR QUANDO DA ERROR
    return false;
  }
}

export function getCharacterXML() {
  // @ts-ignore
  return window.external.getCharacterInfo();
}

export interface Character {
  uid: string;
  name: string;
  weapon: Weapon;
  HR: number;
  GR: number;
  gender: CharacterGender;
  lastLogin: Date;
}

export function getCharacters(): Character[] {
  const parser = new DOMParser();
  const xml = getCharacterXML();
  if (!xml) return [];

  const doc = parser.parseFromString(xml, 'text/xml');

  const docs = doc.getElementsByTagName('Character');
  if (!docs.length) return [];

  const chars: Character[] = [];
  for (const i in docs) {
    const attributes = docs[i].attributes;

    if (!attributes) continue;

    const char = {
      uid: attributes?.uid?.value!,
      name: attributes?.name?.value!,
      weapon: parserWeapon(attributes?.weapon?.value!),
      HR: parseInt(attributes?.HR?.value!, 10),
      GR: parseInt(attributes?.GR?.value!, 10),
      gender: parseCharGender(attributes?.sex?.value!),
      lastLogin: parseCharLastLogin(attributes?.lastLogin?.value!),
    };

    chars.push(char);
  }

  return chars;
}

export enum Weapon {
  SwordAndShield = 'Sword & Shield',
  DualSwords = 'Dual Swords',
  Greatsword = 'Greatsword',
  Longsword = 'Longsword',
  Hammer = 'Hammer',
  HuntingHorn = 'Hunting Horn',
  Lance = 'Lance',
  Gunlance = 'Gunlance',
  Tonfa = 'Tonfa',
  SwitchAxe = 'Switch Axe',
  MagnetSpike = 'Magnet Spike',
  HeavyBowgun = 'Heavy Bowgun',
  LightBowgun = 'Light Bowgun',
  Bow = 'Bow',
  Unknown = 'Unknown',
}

function parserWeapon(weapon: string): Weapon {
  switch (weapon) {
    case '片手剣':
      return Weapon.SwordAndShield;
    case '双剣':
      return Weapon.DualSwords;
    case '大剣':
      return Weapon.Greatsword;
    case '太刀':
      return Weapon.Longsword;
    case 'ハンマー':
      return Weapon.Hammer;
    case '狩猟笛':
      return Weapon.HuntingHorn;
    case 'ランス':
      return Weapon.Lance;
    case 'ガンランス':
      return Weapon.Gunlance;
    case '穿龍棍':
      return Weapon.Tonfa;
    case 'スラッシュアックスＦ':
      return Weapon.SwitchAxe;
    case 'マグネットスパイク':
      return Weapon.MagnetSpike;
    case 'ヘビィボウガン':
      return Weapon.HeavyBowgun;
    case 'ライトボウガン':
      return Weapon.LightBowgun;
    case '弓':
      return Weapon.Bow;
    default:
      return Weapon.Unknown;
  }
}

enum CharacterGender {
  Male = 'Male',
  Famele = 'Famele',
}
function parseCharGender(gender: string): CharacterGender {
  if (gender === 'M') return CharacterGender.Male;
  return CharacterGender.Famele;
}

function parseCharLastLogin(lastLogin: string): Date {
  return new Date(1e3 * parseInt(lastLogin, 10));
}

export function getUserId(): string {
  //@ts-ignore
  return window.external.getUserId();
}

export const isNeAccountChar = (char: Character) => {
  return char.name === '??????' && char.weapon === Weapon.Unknown;
};
