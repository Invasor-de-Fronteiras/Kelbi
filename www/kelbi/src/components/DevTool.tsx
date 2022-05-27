/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/ban-ts-comment */
import React, { useEffect, useMemo, useState } from 'react';
import { randomHexColor } from '../utils/util';
import { BiShow } from 'react-icons/bi';

// data from original launcher
const _UPD_BAR_WID = 302;
const _UPD_BAR_PER = 0.01 * _UPD_BAR_WID;

/**
 * Component for debugging
 */
export function DevTool() {
  const [data, setData] = useState({});
  const [err, setErr] = useState(null);
  const [show, setShow] = useState(true);

  useEffect(() => {
    const timeout = setInterval(() => {
      try {
        const updatePercentageTotal = window.external.getUpdatePercentageTotal();
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
        const updateStatus = window.external.getUpdateStatus();
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
      }}
    >
      <div className='flex flex-row items-center justify-center'>
        <h1>Dev Tool</h1>
        <BiShow style={{ marginLeft: 10 }} size={25} onClick={() => setShow(!show)} />
      </div>
      <h2>States</h2>
      {show && (
        <>
          <div style={{ display: 'flex', flexDirection: 'column' }}>
            {(Object.keys(data) as (keyof typeof data)[])
              .filter((key) => !!data[key])
              .map((key) => {
                return <DebugItem data={data[key]} key={key} name={key} />;
              })}
            {/** @ts-ignore */}
            {data?.updatePercentageTotal && (
              <>
                <span>
                  {/** @ts-ignore */}
                  calc updatePercentageTotal: {Math.ceil(data.updatePercentageTotal * _UPD_BAR_PER)}
                </span>
              </>
            )}
          </div>
          {err && (
            <div style={{ display: 'flex', flexDirection: 'column' }}>
              <h2>Error</h2>
              {/** @ts-ignore */}
              <span>name: {err?.name}</span>
              {/** @ts-ignore */}
              <span>msg: {err?.message}</span>
              {/** @ts-ignore */}
              <span>stack: {err?.stack}</span>
            </div>
          )}
        </>
      )}
    </div>
  );
}

interface DebugItemProps {
  name: string;
  data: any;
}

function DebugItem({ data, name }: DebugItemProps) {
  const [accordion, setAccordion] = useState(false);
  const color = useMemo(() => randomHexColor(), []);
  const value = useMemo(() => (typeof data === 'object' ? JSON.stringify(data) : data), [data]);
  return (
    <div key={name}>
      <BiShow onClick={() => setAccordion(!accordion)} />
      <span style={{ color }}>{name}:</span>
      <span>{value}</span>
      {accordion && (
        <input
          style={{
            border: '1px solid #ccc',
          }}
          value={value}
        />
      )}
    </div>
  );
}
