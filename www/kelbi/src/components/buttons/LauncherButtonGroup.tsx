import React from 'react';
import { AiOutlineClose } from 'react-icons/ai';
import { MdOutlineMinimize } from 'react-icons/md';

import './LauncherButtonGroup.css';

export function LauncherButtonGroup() {
  return (
    <div className='launcher-buttons'>
      <MdOutlineMinimize id='minimize-window-button' className='icon' size={15} />
      <AiOutlineClose id='close-window-button' className='icon' size={15} />
    </div>
  );
}
