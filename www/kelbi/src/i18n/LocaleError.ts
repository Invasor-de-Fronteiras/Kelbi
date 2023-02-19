import type { LocaleKeys } from './TranslateProvider';

export class LocaleError extends Error {
  constructor(key: LocaleKeys) {
    super(key);
  }
}
