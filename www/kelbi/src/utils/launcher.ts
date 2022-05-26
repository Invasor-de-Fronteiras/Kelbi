import { links } from '../constants';
import { findByEnd } from './util';

declare global {
  interface External {
    playSound(song: LauncherSongs): void;
    beginDrag(active: boolean): void;
    openBrowser(url: string): void;
    restartMhf(): void;
    minimizeWindow(): void;
    openMhlConfig(): void;
    closeWindow(): void;
    getUpdatePercentageTotal(): unknown;
    getAccountRights(): string;
    getMhfBootMode(): BoostModeTypes;
    getIniLastServerIndex(): unknown;
    getServerListXml(): string;
    getMhfMutexNumber(): number;
    getUserId(): string;
    getPassword(): string;
    getLastAuthResult(): LastAuthResult;
    getSignResult(): SignResult;
    isEnableSessionId(): unknown;
    getCharacterInfo(): string;
    extractLog(): unknown;
    getUpdateStatus(): number;
    getLauncherReturnCode(): 'NORMAL';
    selectCharacter(charUid: string, charUid1: string): void;
    exitLauncher(): void;
    loginCog(username: string, password: string, confirmPassword: string): void;
  }
}

type LauncherSongs = 'IDR_WAV_SEL' | 'IDR_WAV_OK' | 'IDR_WAV_PRE_LOGIN' | 'IDR_NIKU';
type BoostModeTypes = '_MHF_NORMAL';

export enum LastAuthResult {
  None = 'AUTH_NULL',
  AuthSuccess = 'AUTH_SUCCESS',
  InLoading = 'AUTH_PROGRESS',
  AuthErrorAcc = 'AUTH_ERROR_ACC',
}

export enum SignResult {
  None = 'SIGN_UNKNOWN',
  SignSuccess = 'SIGN_SUCCESS',
  NotMatchPassword = 'SIGN_EPASS',
}

export function openBrowser(url: string) {
  window.external.openBrowser(url);
}

export function openDiscord() {
  openBrowser(links.discordUrl);
}

export function openGithub() {
  openBrowser(links.githubUrl);
}

export function openConfig() {
  window.external.openMhlConfig();
}

export function isEnableSessionId() {
  try {
    return window.external.isEnableSessionId();
  } catch {
    // IGNORANDO ERRO IGUAL LAUNCHER ORIGIONAL
    // TODO: DEBUGAR QUANDO DA ERROR
    return false;
  }
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
  const xml = window.external.getCharacterInfo();
  if (!xml) return [];

  const doc = parser.parseFromString(xml, 'text/xml');

  const docs = doc.getElementsByTagName('Character');
  if (!docs.length) return [];

  const chars: Character[] = [];
  for (const i in docs) {
    const attributes = docs[i].attributes;

    if (!attributes) continue;

    const char = {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      //@ts-ignore
      uid: attributes?.uid?.value,
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      //@ts-ignore
      name: attributes?.name?.value,
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      //@ts-ignore
      weapon: parserWeapon(attributes?.weapon?.value),
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      //@ts-ignore
      HR: parseInt(attributes?.HR?.value, 10),
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      //@ts-ignore
      GR: parseInt(attributes?.GR?.value, 10),
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      //@ts-ignore
      gender: parseCharGender(attributes?.sex?.value),
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      //@ts-ignore
      lastLogin: parseCharLastLogin(attributes?.lastLogin?.value),
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

export const isNeAccountChar = (char: Character) => {
  return char.name === '??????' && char.weapon === Weapon.Unknown;
};

export function startGame(charId: string) {
  window.external.selectCharacter(charId, charId);

  window.external.exitLauncher();
}

export function createNewCharacter(): Promise<Character> {
  return new Promise((resolve, reject) => {
    const hasNewChar = (): boolean => {
      const lastAuth = window.external.getLastAuthResult();

      if (lastAuth === LastAuthResult.InLoading) {
        return false;
      }

      const chars = getCharacters();
      const newChar = findByEnd(chars, isNeAccountChar);

      if (!newChar) {
        return false;
      }

      resolve(newChar);
      return true;
    };

    try {
      // checking if the user has a new char that was not used before
      if (hasNewChar()) {
        return;
      }

      const userId = window.external.getUserId();
      const password = window.external.getPassword();

      //creating new char
      window.external.loginCog(userId + '+', password, password);

      const interval = setInterval(() => {
        if (hasNewChar()) {
          clearInterval(interval);
        }
      }, 100);
    } catch (error) {
      reject(error);
    }
  });
}
