import React from 'react';
import { SignInForm } from './components/SignInForm';

import { AiFillGithub } from 'react-icons/ai';
import { GrConfigure } from 'react-icons/gr';
import { SiDiscord } from 'react-icons/si';
import { LogoIcon } from './components/LogoIcon';

import './styles/app.css';
import { Background } from './components/Background';

function App() {
  return (
    <div id='main'>
      <div id='left-side'>
        <div id='left-side-header'>
          <LogoIcon />
        </div>
        <SignInForm />
        <div id='left-side-footer'>
          <GrConfigure size={50} />
          <SiDiscord size={50} />
          <AiFillGithub size={50} />
        </div>
      </div>
      <Background />
    </div>
  );
}

export default App;
