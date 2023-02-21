import { useState } from 'react';
import { UpdateStatus } from '../utils/launcher';

interface LauncherUpdateHook {
  loading: boolean;
  updateOk: boolean;
  version: string | null;
  fileProgress: number;
  totalProgress: number;
  startUpdate: () => void;
  verifyUpdates: () => void;
}

const getCurrentVersion = () => {
  try {
    const xmlHttp = new XMLHttpRequest();
    xmlHttp.open('GET', 'http://cog-members.mhf-z.jp/version', false);
    xmlHttp.send(null);
    if (xmlHttp.status === 404) {
      return null;
    }

    return xmlHttp.responseText;
  } catch (err) {
    return null;
  }
};

export function useLauncherUpdate(): LauncherUpdateHook {
  const [state, setState] = useState<Omit<LauncherUpdateHook, 'startUpdate' | 'verifyUpdates'>>({
    loading: false,
    updateOk: false,
    version: localStorage.getItem('files_version') || null,
    fileProgress: 0,
    totalProgress: 0,
  });

  const startUpdate = () => {
    window.external.startUpdate();

    setState({
      ...state,
      loading: true,
    });

    const interval = setInterval(() => {
      const totalProgress = window.external.getUpdatePercentageTotal();
      const fileProgress = window.external.getUpdatePercentageFile();
      const updateStatus = window.external.getUpdateStatus();

      if (updateStatus === UpdateStatus.None) {
        return;
      }

      if (updateStatus === UpdateStatus.UpdateStart) {
        setState({
          ...state,
          totalProgress,
          fileProgress,
          loading: true,
        });
        return;
      }

      clearInterval(interval);

      if (updateStatus === UpdateStatus.UpdateOk) {
        setState({
          ...state,
          updateOk: true,
          loading: false,
        });
        return;
      }

      setState({
        ...state,
        loading: false,
      });
    }, 100);

    return () => clearInterval(interval);
  };

  const verifyUpdates = () => {
    const currentVersion = getCurrentVersion();

    if (!currentVersion || (state.version && state.version == currentVersion)) {
      setState({
        ...state,
        updateOk: true,
      });
      return;
    }

    localStorage.setItem('files_version', currentVersion);
    startUpdate();
  };

  return {
    ...state,
    startUpdate,
    verifyUpdates,
  };
}
