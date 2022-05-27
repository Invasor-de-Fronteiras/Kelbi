import React, { createContext, useContext, useState } from 'react';
import { Loading } from '../components/Loading';

interface LauncherContextProps {
  isLoading: boolean;
  loggedIn: boolean;
  setIsLoading: (isLoading: boolean) => void;
  setLoggedIn: (loggedIn: boolean) => void;
}

const LauncherContext = createContext({} as LauncherContextProps);

export function LauncherProvider({ children }: { children: React.ReactNode }) {
  const [isLoading, setIsLoading] = useState(false);
  const [loggedIn, setLoggedIn] = useState(false);

  return (
    <LauncherContext.Provider
      value={{
        isLoading,
        loggedIn,
        setLoggedIn,
        setIsLoading,
      }}
    >
      {children}
      {isLoading && <Loading />}
    </LauncherContext.Provider>
  );
}

export const useLauncher = () => useContext(LauncherContext);
