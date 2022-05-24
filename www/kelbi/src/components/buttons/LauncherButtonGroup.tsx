import React from 'react';
import { AiOutlineClose } from 'react-icons/ai';
import { MdOutlineMinimize } from 'react-icons/md';
import { FaQuestion } from 'react-icons/fa';
import { closeWindow, minimizeWindow } from '../../utils/launcher';

import './LauncherButtonGroup.css';

export function LauncherButtonGroup() {
  return (
    <div className='launcher-buttons'>
      <FaQuestion size={15} className='icon' />
      <MdOutlineMinimize id='minimize-window-button' size={15} onClick={minimizeWindow} />
      <AiOutlineClose id='close-window-button' className='icon' size={15} onClick={closeWindow} />
    </div>
  );
}
