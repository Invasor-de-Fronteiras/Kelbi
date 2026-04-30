import React, { useEffect, useState } from 'react';
import { Button } from '../components/Button';
import { useLauncher } from '../context/LauncherContext';
import { useTranslate } from '../i18n/useTranslate';
import { Server, getServerList } from '../utils/launcher';

interface SelectServerProps {
  onClose?: () => void;
}

export function SelectServer({ onClose }: SelectServerProps) {
  const { t } = useTranslate();
  const { setSelectedServer, selectedServer } = useLauncher();
  const [servers, setServers] = useState<Server[]>([]);
  const [selectedIdx, setSelectedIdx] = useState<number | null>(null);

  useEffect(() => {
    const list = getServerList();
    setServers(list);

    if (list.length === 1 && !onClose) {
      window.external.setIniLastServerIndex(list[0].idx);
      setSelectedServer(list[0]);
      return;
    }

    if (list.length > 0) {
      const preferredIdx =
        onClose && selectedServer != null
          ? selectedServer.idx
          : window.external.getIniLastServerIndex();
      const preferred = list.find((s) => s.idx === preferredIdx);
      setSelectedIdx(preferred ? preferred.idx : list[0].idx);
    }
  }, []);

  const handleEnter = () => {
    const server = servers.find((s) => s.idx === selectedIdx);
    if (!server) return;
    window.external.setIniLastServerIndex(server.idx);
    setSelectedServer(server);
    if (onClose) onClose();
  };

  return (
    <div className='flex flex-col items-center'>
      <h3>{t('select_server_title')}</h3>
      <div style={{ marginBottom: '16px' }}>
        {servers.map((server) => (
          <div
            key={server.idx}
            className='character-card'
            role='option'
            onClick={() => setSelectedIdx(server.idx)}
            style={{
              backgroundColor: server.idx === selectedIdx ? '#f83' : '#fff',
              cursor: 'pointer',
            }}
          >
            <h3 className='char-name'>{server.name}</h3>
            <span className='char-weapon'>
              {server.ip}:{server.port}
            </span>
          </div>
        ))}
      </div>
      <Button onClick={handleEnter} disabled={selectedIdx === null}>
        {t('select_server_enter')}
      </Button>
      {onClose && (
        <Button onClick={onClose}>{t('select_server_cancel')}</Button>
      )}
    </div>
  );
}
