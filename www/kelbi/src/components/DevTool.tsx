import React from 'react';
import { useLauncher } from '../context/LauncherContext';

export function DevTool() {
  const { error, lastAuthResult, signResult, enableSessionId, charset } = useLauncher();

  return (
    <div
      style={{
        backgroundColor: '#fff',
        position: 'absolute',
        zIndex: 1,
        width: 'min-content',
        maxWidth: '50%',
        height: 'min-content',
        top: 0,
        right: 0,
        marginRight: '10%',
      }}>
      <h1>Dev Tool</h1>
      <div style={{ display: 'flex', flexDirection: 'column' }}>
        <h2>States</h2>
        <span>lastAuthResult: {lastAuthResult}</span>
        <span>signResult: {signResult}</span>
        <span>enableSessionId: {enableSessionId}</span>
        <span>charset: {charset}</span>
      </div>
      {error && (
        <div style={{ display: 'flex', flexDirection: 'column' }}>
          <h2>Error</h2>
          <span>name: {error.name}</span>
          <span>msg: {error.message}</span>
          <span>stack: {error.stack}</span>
        </div>
      )}
    </div>
  );
}
