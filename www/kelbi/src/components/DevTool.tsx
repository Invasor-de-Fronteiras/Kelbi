/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/ban-ts-comment */
import React, { useEffect, useMemo, useState } from 'react';

import './DevTool.css';

// data from original launcher
const _UPD_BAR_WID = 302;
const _UPD_BAR_PER = 0.01 * _UPD_BAR_WID;

/**
 * Component for debugging
 */
export function DevTool() {
  const [data, setData] = useState({});
  const [showLogin, setShowLogin] = useState(false);
  const [startUpdateResult, setStartUpdateResult] = useState(null);

  useEffect(() => {
    const timeout = setInterval(() => {
      try {
        const accountRights = window.external.getAccountRights();
        const mhfBootMode = window.external.getMhfBootMode();
        const lastServerIndex = window.external.getIniLastServerIndex();
        const serverListXML = window.external.getServerListXml();
        const mhfMutexNumber = window.external.getMhfMutexNumber();
        const userId = window.external.getUserId();
        const password = window.external.getPassword();
        const lastAuthResult = window.external.getLastAuthResult();
        const signResult = window.external.getSignResult();
        const enableSessionId = window.external.isEnableSessionId();
        const charsXML = window.external.getCharacterInfo();
        const extractLog = window.external.extractLog();
        const launcherReturnCode = window.external.getLauncherReturnCode();
        const updateStatus = window.external.getUpdateStatus();
        const updatePercentageTotal = window.external.getUpdatePercentageTotal();
        const updatePercentageFile = window.external.getUpdatePercentageFile();

        setData({
          browserLang: navigator.language,
          accountRights,
          mhfBootMode,
          lastServerIndex,
          serverListXML,
          mhfMutexNumber,
          lastAuthResult,
          signResult,
          enableSessionId,
          charsXML,
          extractLog,
          launcherReturnCode,
          startUpdateResult,
          updateStatus,
          updatePercentageFile,
          updatePercentageTotal,
          calcUpdatePercentageTotal: Math.ceil(Number(updatePercentageTotal) * _UPD_BAR_PER),
          ...(showLogin ? { userId, password } : {}),
        });
      } catch (err) {
        setData((prev) => ({
          ...prev,
          errName: err?.name,
          errMsg: err?.message,
          errStack: err?.stack,
          errTime: new Date(),
        }));
      }
    }, 100);
    return () => clearInterval(timeout);
  }, [showLogin, startUpdateResult]);

  const onStartUpdate = () => {
    const result = window.external.startUpdate();
    setStartUpdateResult(result);
  };

  return (
    <div id='debugger'>
      <div className='header'>
        <h1>Dev Tool</h1>
        <p>
          Painel de Debug do Launcher, para entender seus efeitos colateiras e acompanhar os
          estados. Usado para encontrar erros no launcher.
        </p>
        <div id='show-login'>
          <input type='checkbox' checked={showLogin} onClick={() => setShowLogin((e) => !e)} />
          <label>Show Login</label>
        </div>
        <button onClick={onStartUpdate}>Start update</button>
      </div>
      <table>
        <caption>WATCH</caption>
        <thead>
          <tr>
            <th>State</th>
            <th>Type</th>
            <th>Value</th>
            <th>Copy</th>
          </tr>
        </thead>
        <tbody>
          {(Object.keys(data) as (keyof typeof data)[])
            .filter((key) => !!data[key])
            .map((key) => {
              return <DebugItem data={data[key]} key={key} name={key} />;
            })}
        </tbody>
      </table>
    </div>
  );
}

interface DebugItemProps {
  name: string;
  data: any;
}

function DebugItem({ data, name }: DebugItemProps) {
  const value = useMemo(
    () => (typeof data === 'object' ? JSON.stringify(data) : data?.toString() ?? data),
    [data],
  );

  return (
    <tr key={name}>
      <td>
        <strong>{name}</strong>
      </td>
      <td>{typeof data}</td>
      <td>{value}</td>
      <td>
        <input value={value} />
      </td>
    </tr>
  );
}
