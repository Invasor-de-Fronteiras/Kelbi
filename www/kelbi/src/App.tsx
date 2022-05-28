import React from 'react';
import { SignInForm } from './layout/SignInForm';

import { useLauncher } from './context/LauncherContext';
import { Layout } from './layout/Layout';
import { SelectCharacter } from './layout/SelectCharacter';

function App() {
  const { loggedIn } = useLauncher();

  return <Layout>{loggedIn ? <SelectCharacter /> : <SignInForm />}</Layout>;
}

export default App;
