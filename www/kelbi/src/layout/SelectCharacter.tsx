import React, { useEffect, useState } from 'react';

import { Button } from '../components/Button';
import { CharacterCard } from '../components/CharacterCard';
import { LauncherLoading } from '../components/LauncherLoading';
import { useCreateCharacter } from '../hooks/useCreateCharacter';
import { useGetCharacters } from '../hooks/useGetCharacters';
import { useLauncherUpdate } from '../hooks/useLauncherUpdate';
import { LocaleKeys } from '../i18n/TranslateProvider';
import { useTranslate } from '../i18n/useTranslate';
import { startGame } from '../utils/launcher';
import { playLoginSong, randomSong } from '../utils/songs';

export function SelectCharacter() {
  const { t } = useTranslate();

  const { mutate: handleCreateNewChar, isLoading: newCharInLoading } = useCreateCharacter();
  const [localLoadingMessage, setLoadingMessage] = useState<LocaleKeys | null>(null);
  const {
    startUpdate,
    updateOk,
    fileProgress,
    totalProgress,
    loading: updateLoading,
  } = useLauncherUpdate();
  const { characters, isNewAccount, newAccountUID, loading: charLoading } = useGetCharacters();
  const [selectedCharId, setSelectedCharId] = useState('');

  useEffect(() => {
    if (characters[0]) {
      setSelectedCharId(characters[0].uid);
    }
  }, [characters]);

  useEffect(() => {
    if (!charLoading) {
      startUpdate();
    }
  }, [charLoading]);

  useEffect(() => {
    if (updateLoading) {
      setLoadingMessage('update_files_loading');
    } else {
      setLoadingMessage(null);
    }
  }, [updateLoading]);

  const loading = localLoadingMessage !== null || charLoading || !updateOk;
  const loadingMessage = localLoadingMessage ? t(localLoadingMessage) : t('searching_data_loading');

  const handleStartGame = () => {
    playLoginSong();
    setLoadingMessage('opening_game_loading');
    startGame(selectedCharId ?? newAccountUID);
  };

  const handleChangeAccount = () => {
    setLoadingMessage('change_account_loading');

    localStorage.removeItem('autoLogin');
    localStorage.removeItem('password');
    localStorage.removeItem('username');
    window.external.restartMhf();
  };

  if (loading) {
    return (
      <LauncherLoading
        message={loadingMessage}
        fileProgress={fileProgress}
        totalProgress={totalProgress}
        showUpdateProgress={updateLoading}
      />
    );
  }

  if (isNewAccount) {
    return (
      <div>
        <h4 className='text-center'>{t('new_account_message')}</h4>
        <div className='flex items-center justify-center flex-col'>
          <Button onClick={handleStartGame}>{t('new_account_enter')}</Button>
          <Button onClick={handleChangeAccount}>{t('new_account_change_account')}</Button>
        </div>
      </div>
    );
  }

  const layoutWidth = 300;

  return (
    <div>
      <div
        style={{
          overflowY: 'auto',
          maxHeight: '70%',
          paddingRight: 10,
          marginTop: 50,
          top: 0,
          position: 'absolute',
          marginLeft: layoutWidth + 100,
        }}
      >
        {characters.map((char, index) => (
          <CharacterCard
            char={char}
            key={char.uid}
            tabIndex={index}
            isSelected={char.uid === selectedCharId}
            onSelect={() => {
              randomSong();
              setSelectedCharId((id) => (id === char.uid ? null : char.uid));
            }}
          />
        ))}
      </div>
      <div className='flex flex-col items-center'>
        <h3 style={{}}>{t('select_char_title')}</h3>
        <Button onClick={handleStartGame} disabled={!selectedCharId || newCharInLoading}>
          {t('select_char_enter')}
        </Button>
        <Button
          onClick={handleCreateNewChar}
          loadingMessage={t('new_char_loading')}
          isLoading={newCharInLoading}
        >
          {t('select_char_new_char')}
        </Button>
        <Button onClick={handleChangeAccount} disabled={newCharInLoading}>
          {t('select_char_change_account')}
        </Button>
        <AutoLoginCheckbox />
      </div>
    </div>
  );
}

function AutoLoginCheckbox() {
  const [checked, setChecked] = useState(() => localStorage.getItem('autoLogin') === 'true');
  const { t } = useTranslate();

  function handleChange(checked: boolean) {
    localStorage.setItem('autoLogin', String(checked));
    setChecked(checked);
  }

  return (
    <div
      className='checkbox-group'
      style={{
        marginTop: '10px',
      }}
      onClick={() => handleChange(!checked)}
    >
      <input
        type='checkbox'
        id='auto-login'
        onChange={(e) => handleChange(e.target.checked)}
        checked={checked}
      />
      <label htmlFor='auto-login'>{t('select_char_keep_account')}</label>
    </div>
  );
}
