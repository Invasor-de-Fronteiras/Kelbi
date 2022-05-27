import React from 'react';

import { AiFillGithub } from 'react-icons/ai';
import { GrConfigure } from 'react-icons/gr';
import { SiDiscord } from 'react-icons/si';
import { LogoIcon } from '../components/LogoIcon';

import { Background } from '../components/Background';
import { LauncherButtonGroup } from '../components/buttons/LauncherButtonGroup';
import { openConfig, openDiscord, openGithub } from '../utils/launcher';

import './Layout.css';
import { DevTool } from '../components/DevTool';

interface LayoutProps {
  children: React.ReactNode;
}

export function Layout({ children }: LayoutProps) {
  return (
    <>
      <div id='main'>
        <div id='left-side'>
          <div id='left-side-header'>
            <LogoIcon />
          </div>
          <div id='left-side-body'>{children}</div>
          <div id='left-side-footer'>
            <GrConfigure className='icon' size={20} onClick={openConfig} />
            <div>
              <SiDiscord className='icon discord-icon' size={20} onClick={openDiscord} />
              <AiFillGithub className='icon' size={20} onClick={openGithub} />
            </div>
          </div>
        </div>
      </div>
      {/* <DevTool /> */}
      <Background />
      <LauncherButtonGroup />
    </>
  );
}
