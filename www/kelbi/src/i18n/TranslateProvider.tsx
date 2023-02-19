import React, { createContext, useCallback, useState } from 'react';
import { LocaleError } from './LocaleError';
import { enUSLocales } from './locales/en-us';
import { esESLocales } from './locales/es-ES';
import { ptBRLocales } from './locales/pt-br';

type Locale = 'pt-BR' | 'ja' | 'es-ES' | 'en-US';
export type LocaleKeys = keyof typeof ptBRLocales;

type TFunction = (key: LocaleKeys | LocaleError) => string;
export const TranslateContext = createContext(
  {} as {
    locale: Locale;
    setLocale: (newLang: Locale) => unknown;
    t: TFunction;
  },
);

const locales: Record<Locale, Record<LocaleKeys, string>> = {
  'pt-BR': ptBRLocales,
  'en-US': enUSLocales,
  'es-ES': esESLocales,
};

export const TranslateProvider = ({ children }: { children: React.ReactNode }) => {
  const [locale, setLocale] = useState<Locale>(() => {
    const storedLocale = localStorage.getItem('language') as Locale | null;
    return storedLocale ?? 'pt-BR';
  });

  const t: TFunction = useCallback(
    (k) => {
      const key = k instanceof LocaleError ? k.message : k;
      return locales[locale]?.[key] ?? key;
    },
    [locale],
  );

  const setLocalePersistent = (locale: Locale) => {
    localStorage.setItem('language', locale);
    setLocale(locale);
  };
  return (
    <TranslateContext.Provider value={{ locale, setLocale: setLocalePersistent, t }}>
      {children}
    </TranslateContext.Provider>
  );
};
