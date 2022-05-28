import React from 'react';
import { AiOutlineClose } from 'react-icons/ai';
import { MdOutlineMinimize } from 'react-icons/md';

import './LauncherButtonGroup.css';

export function LauncherButtonGroup() {
  return (
    <div className='launcher-buttons'>
      {/* <FaQuestion size={15} className='icon' /> */}
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
