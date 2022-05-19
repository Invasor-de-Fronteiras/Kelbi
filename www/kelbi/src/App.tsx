import React from 'react';
import { SignInForm } from './layout/SignInForm';

import { useLauncher } from './context/LauncherContext';
import { Layout } from './layout/Layout';

function App() {
  const { loggedIn } = useLauncher();

  return <Layout>{loggedIn ? <span>VC ESTÁ LOGADO</span> : <SignInForm />}</Layout>;
}

export default App;
