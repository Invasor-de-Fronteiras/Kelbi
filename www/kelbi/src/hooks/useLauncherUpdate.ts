import { useState } from 'react';
import { UpdateStatus } from '../utils/launcher';

interface LauncherUpdateHook {
  loading: boolean;
  updateOk: boolean;
  fileProgress: number;
  totalProgress: number;
  startUpdate: () => void;
}

export function useLauncherUpdate(): LauncherUpdateHook {
  const [state, setState] = useState<Omit<LauncherUpdateHook, 'startUpdate'>>({
    loading: false,
    updateOk: false,
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

  return {
    ...state,
    startUpdate,
  };
}
