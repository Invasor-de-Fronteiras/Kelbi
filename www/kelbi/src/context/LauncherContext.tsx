import React, { createContext, useContext, useState } from 'react';
import { Loading } from '../components/Loading';
import type { Server } from '../utils/launcher';

interface LauncherContextProps {
  showDebugger: boolean;
  setShowDebugger: (active: boolean) => void;
  loggedIn: boolean;
  isLoading: boolean;
  setIsLoading: (isLoading: boolean) => void;
  setLoggedIn: (loggedIn: boolean) => void;
  selectedServer: Server | null;
  setSelectedServer: (server: Server) => void;
}

const LauncherContext = createContext({} as LauncherContextProps);

export function LauncherProvider({ children }: { children: React.ReactNode }) {
  const [isLoading, setIsLoading] = useState(false);
  const [loggedIn, setLoggedIn] = useState(false);
  const [showDebugger, setShowDebugger] = useState(false);
  const [selectedServer, setSelectedServer] = useState<Server | null>(null);

  return (
    <LauncherContext.Provider
      value={{
        showDebugger,
        setShowDebugger,
        isLoading,
        loggedIn,
        setLoggedIn,
        setIsLoading,
        selectedServer,
        setSelectedServer,
      }}
    >
      {children}
      {isLoading && <Loading />}
    </LauncherContext.Provider>
  );
}

export const useLauncher = () => useContext(LauncherContext);
