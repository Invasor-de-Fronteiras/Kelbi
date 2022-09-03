import React, { createContext, useContext, useState } from 'react';
import { Loading } from '../components/Loading';

interface LauncherContextProps {
  showDebugger: boolean;
  setShowDebugger: (active: boolean) => void;
  loggedIn: boolean;
  isLoading: boolean;
  setIsLoading: (isLoading: boolean) => void;
  setLoggedIn: (loggedIn: boolean) => void;
}

const LauncherContext = createContext({} as LauncherContextProps);

export function LauncherProvider({ children }: { children: React.ReactNode }) {
  const [isLoading, setIsLoading] = useState(false);
  const [loggedIn, setLoggedIn] = useState(false);
  const [showDebugger, setShowDebugger] = useState(false);

  return (
    <LauncherContext.Provider
      value={{
        showDebugger,
        setShowDebugger,
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
