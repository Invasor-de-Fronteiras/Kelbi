import React, { createContext, useContext, useState } from 'react';
import { Loading } from '../components/Loading';

interface LauncherContextProps {
  isLoading: boolean;
  loggedIn: boolean;
  setLoggedIn: (loggedIn: boolean) => void;
  setIsLoading: (isLoading: boolean) => void;
}

const LauncherContext = createContext({} as LauncherContextProps);

export function LauncherProvider({ children }: any) {
  const [isLoading, setIsLoading] = useState(false);
  const [loggedIn, setLoggedIn] = useState(false);

  return (
    <LauncherContext.Provider value={{ isLoading, setIsLoading, loggedIn, setLoggedIn }}>
      {children}
      {isLoading && <Loading />}
    </LauncherContext.Provider>
  );
}

export const useLauncher = () => useContext(LauncherContext);
