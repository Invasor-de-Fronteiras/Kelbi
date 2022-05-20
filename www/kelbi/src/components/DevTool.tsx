import React, { useEffect, useState } from 'react';
import { randomHexColor } from '../utils/util';

// data from original launcher
const _UPD_BAR_WID = 302;
const _UPD_BAR_PER = 0.01 * _UPD_BAR_WID;

export function DevTool() {
  const [data, setData] = useState({
    updatePercentageTotal: null,
    accountRights: null,
    mhfBootMode: null,
    lastServerIndex: null,
    serverListXML: null,
    mhfMutexNumber: null,
    userId: null,
    password: null,
    lastAuthResult: null,
    signResult: null,
    enableSessionId: null,
    charsXML: null,
    extractLog: null,
    updateStatus: null,
    launcherReturnCode: null,
  });
  const [err, setErr] = useState(null);

  useEffect(() => {
    let timeout = setInterval(() => {
      try {
        //@ts-ignore
        const updatePercentageTotal = window.external.getUpdatePercentageTotal();
        //@ts-ignore
        const accountRights = window.external.getAccountRights();
        //@ts-ignore
        const mhfBootMode = window.external.getMhfBootMode();
        //@ts-ignore
        const lastServerIndex = window.external.getIniLastServerIndex();
        //@ts-ignore
        const serverListXML = window.external.getServerListXml();
        //@ts-ignore
        const mhfMutexNumber = window.external.getMhfMutexNumber();
        //@ts-ignore
        const userId = window.external.getUserId();
        //@ts-ignore
        const password = window.external.getPassword();
        //@ts-ignore
        const lastAuthResult = window.external.getLastAuthResult();
        //@ts-ignore
        const signResult = window.external.getSignResult();
        //@ts-ignore
        const enableSessionId = window.external.isEnableSessionId();
        //@ts-ignore
        const charsXML = window.external.getCharacterInfo();
        //@ts-ignore
        const extractLog = window.external.extractLog();
        //@ts-ignore
        const updateStatus = window.external.getUpdateStatus();
        //@ts-ignore
        const launcherReturnCode = window.external.getLauncherReturnCode();

        setData({
          updatePercentageTotal,
          accountRights,
          mhfBootMode,
          lastServerIndex,
          serverListXML,
          mhfMutexNumber,
          userId,
          password,
          lastAuthResult,
          signResult,
          enableSessionId,
          charsXML,
          extractLog,
          updateStatus,
          launcherReturnCode,
        });
      } catch (err) {
        //@ts-ignore
        setErr(err);
      }
    }, 100);
    return () => clearInterval(timeout);
  }, []);

  return (
    <div
      style={{
        backgroundColor: '#fff',
        position: 'absolute',
        zIndex: 1,
        padding: '5px',
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
        {(Object.keys(data) as (keyof typeof data)[])
          .filter((key) => !!data[key])
          .map((key) => {
            return (
              <div key={key}>
                <span style={{ color: randomHexColor() }}>{key}:</span>
                <span>{data[key]}</span>
              </div>
            );
          })}
        {data?.updatePercentageTotal && (
          <>
            <span>
              calc updatePercentageTotal: {Math.ceil(data.updatePercentageTotal * _UPD_BAR_PER)}
            </span>
          </>
        )}
      </div>
      {err && (
        <div style={{ display: 'flex', flexDirection: 'column' }}>
          <h2>Error</h2>
          <span>name: {err?.name}</span>
          <span>msg: {err?.message}</span>
          <span>stack: {err?.stack}</span>
        </div>
      )}
    </div>
  );
}
