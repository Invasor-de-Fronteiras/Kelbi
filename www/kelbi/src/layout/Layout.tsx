import React, { useState } from 'react';

import { AiFillGithub } from 'react-icons/ai';
import { GrConfigure } from 'react-icons/gr';
import { SiDiscord } from 'react-icons/si';
import { LogoIcon } from '../components/LogoIcon';

import { Background } from '../components/Background';
import { LauncherButtonGroup } from '../components/buttons/LauncherButtonGroup';
import { openConfig, openDiscord, openGithub } from '../utils/launcher';

import './Layout.css';
import { DevTool } from '../components/DevTool';
import { useLauncher } from '../context/LauncherContext';
import { ChangeLocaleButton } from '../i18n/ChangeLocaleButton';
import { SelectServer } from './SelectServer';

interface LayoutProps {
  children: React.ReactNode;
}

export const TransparentBackground = () => (
  <div
    style={{
      backgroundColor: '#f8f8f8',
      width: '100%',
      height: '100%',
      position: 'absolute',
      opacity: 0.8,
      zIndex: -999,
      top: 0,
      left: 0,
    }}
  />
);

export function Layout({ children }: LayoutProps) {
  const { showDebugger, selectedServer } = useLauncher();
  const [showPicker, setShowPicker] = useState(false);

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
            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
              <SiDiscord className='icon discord-icon' size={20} onClick={openDiscord} />
              <AiFillGithub className='icon' size={20} onClick={openGithub} />
              <ChangeLocaleButton />
            </div>
          </div>
          <TransparentBackground />
        </div>
      </div>

      {selectedServer && (
        <>
          <div
            onClick={() => setShowPicker(true)}
            style={{
              position: 'fixed',
              bottom: 12,
              right: 12,
              backgroundColor: 'rgba(0,0,0,0.65)',
              color: '#fff',
              padding: '4px 10px',
              borderRadius: 4,
              cursor: 'pointer',
              fontSize: '0.8em',
              letterSpacing: '0.05em',
              zIndex: 100,
            }}
          >
            {selectedServer.name}
          </div>

          {showPicker && (
            <div
              style={{
                position: 'fixed',
                top: 0,
                left: 0,
                width: '100%',
                height: '100%',
                backgroundColor: 'rgba(0,0,0,0.55)',
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                zIndex: 200,
              }}
            >
              <div
                style={{
                  backgroundColor: '#f8f8f8',
                  padding: '24px 32px',
                  borderRadius: 4,
                  opacity: 1,
                }}
              >
                <SelectServer onClose={() => setShowPicker(false)} />
              </div>
            </div>
          )}
        </>
      )}

      {showDebugger ? <DevTool /> : null}
      <Background />
      <LauncherButtonGroup />
    </>
  );
}
