import React from 'react';
import { useLauncher } from './context/LauncherContext';
import { Layout } from './layout/Layout';
import { SelectCharacter } from './layout/SelectCharacter';
import { SelectServer } from './layout/SelectServer';
import { SignInForm } from './layout/SignInForm';

function App() {
  const { loggedIn, selectedServer } = useLauncher();

  if (!loggedIn) return <Layout><SignInForm /></Layout>;
  if (!selectedServer) return <Layout><SelectServer /></Layout>;
  return <Layout><SelectCharacter /></Layout>;
}

export default App;
