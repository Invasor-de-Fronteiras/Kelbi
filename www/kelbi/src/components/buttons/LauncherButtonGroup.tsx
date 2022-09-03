import React from 'react';
import { AiOutlineClose } from 'react-icons/ai';
import { MdOutlineMinimize } from 'react-icons/md';
import { VscDebug } from 'react-icons/vsc';
import { useLauncher } from '../../context/LauncherContext';

import './LauncherButtonGroup.css';

export function LauncherButtonGroup() {
  const { setShowDebugger } = useLauncher();

  return (
    <div className='launcher-buttons'>
      <VscDebug size={15} className='icon' onClick={() => setShowDebugger((v) => !v)} />
      <MdOutlineMinimize
        id='minimize-window-button'
        size={15}
        onClick={() => window.external.minimizeWindow()}
      />
      <AiOutlineClose
        id='close-window-button'
        className='icon'
        size={15}
        onClick={() => window.external.closeWindow()}
      />
    </div>
  );
}
