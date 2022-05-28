import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import { LauncherProvider } from './context/LauncherContext';

import './styles/main.css';
import './styles/framework.css';
import './utils/faker-launcher';

ReactDOM.render(
  <React.StrictMode>
    <LauncherProvider>
      <App />
    </LauncherProvider>
  </React.StrictMode>,
  document.getElementById('root'),
);
