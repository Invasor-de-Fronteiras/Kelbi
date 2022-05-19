import React, { createContext, useContext, useEffect, useState } from 'react';
import { Loading } from '../components/Loading';
import {
  getCharacterXML,
  getLastAuthResult,
  getSignResult,
  isEnableSessionId,
  LastAuthResult,
  SignResult,
} from '../utils/launcher';

interface LauncherContextProps {
  isLoading: boolean;
  loggedIn: boolean;
  enableSessionId: boolean;
  lastAuthResult: LastAuthResult;
  signResult: SignResult;
  charset: string;
  error: Error | null;
  setIsLoading: (isLoading: boolean) => void;
}

const LauncherContext = createContext({} as LauncherContextProps);

export function LauncherProvider({ children }: any) {
  const [isLoading, setIsLoading] = useState(false);
  const [loggedIn, setLoggedIn] = useState(false);
  const [lastAuthResult, setLastAuthResult] = useState(LastAuthResult.None);
  const [signResult, setSignResult] = useState(SignResult.None);
  const [charset, setCharset] = useState('');
  const [error, setError] = useState(null);
  const [enableSessionId, setIsEnableSessionId] = useState(false);

  useEffect(() => {
    let timeout = setInterval(() => {
      try {
        const lastAuth = getLastAuthResult();
        const signRes = getSignResult();
        const enableSessionId = isEnableSessionId();
        const charsXML = getCharacterXML();
        //@ts-ignore
        const extracLog = window.external.extractLog();

        if (lastAuth === LastAuthResult.InLoading) {
          setIsLoading(true);
        } else {
          setIsLoading(false);
        }

        if (lastAuth === LastAuthResult.AuthSuccess && signRes === SignResult.SignSuccess) {
          setLoggedIn(true);
        }

        setIsEnableSessionId(enableSessionId);
        setLastAuthResult(lastAuth);
        setSignResult(signRes);
        setCharset(JSON.stringify(charsXML) + String(extracLog));
      } catch (err) {
        //@ts-ignore
        setError(err);
      }
    }, 100);

    return () => {
      clearInterval(timeout);
    };
  }, []);

  return (
    <LauncherContext.Provider
      value={{
        isLoading,
        loggedIn,
        error,
        lastAuthResult,
        charset,
        signResult,
        setIsLoading,
        enableSessionId,
      }}>
      {children}
      {isLoading && <Loading />}
    </LauncherContext.Provider>
  );
}

export const useLauncher = () => useContext(LauncherContext);
