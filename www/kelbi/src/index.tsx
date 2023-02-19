import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import { LauncherProvider } from './context/LauncherContext';

import './styles/main.css';
import './styles/framework.css';
import './utils/faker-launcher';
import { TranslateProvider } from './i18n/TranslateProvider';

ReactDOM.render(
  <React.StrictMode>
    <TranslateProvider>
      <LauncherProvider>
        <App />
      </LauncherProvider>
    </TranslateProvider>
  </React.StrictMode>,
  document.getElementById('root'),
);
