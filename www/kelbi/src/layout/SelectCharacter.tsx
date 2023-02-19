import React, { useEffect, useState } from 'react';

import { Button } from '../components/Button';
import { CharacterCard } from '../components/CharacterCard';
import { LauncherLoading } from '../components/LauncherLoading';
import { useCreateCharacter } from '../hooks/useCreateCharacter';
import { useGetCharacters } from '../hooks/useGetCharacters';
import { useLauncherUpdate } from '../hooks/useLauncherUpdate';
import { openDiscord, startGame } from '../utils/launcher';
import { playLoginSong, randomSong } from '../utils/songs';

export function SelectCharacter() {
  const { mutate: handleCreateNewChar, isLoading: newCharInLoading } = useCreateCharacter();
  const [localLoadingMessage, setLoadingMessage] = useState<string | null>(null);
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
      setLoadingMessage('Atualizando arquivos!');
    } else {
      setLoadingMessage(null);
    }
  }, [updateLoading]);

  const loading = localLoadingMessage !== null || charLoading || !updateOk;
  const loadingMessage = localLoadingMessage ?? 'Buscando seus dados...';

  const handleStartGame = () => {
    playLoginSong();
    setLoadingMessage('Abrindo o jogo...');
    startGame(selectedCharId ?? newAccountUID);
  };

  const handleChangeAccount = () => {
    setLoadingMessage('Trocando de conta...');

    localStorage.removeItem('autoLogin');
    localStorage.removeItem('password');
    localStorage.removeItem('username');
    window.external.restartMhf();
  };

  if (loading)
    return (
      <LauncherLoading
        message={loadingMessage}
        fileProgress={fileProgress}
        totalProgress={totalProgress}
        showUpdateProgress={updateLoading}
      />
    );

  if (isNewAccount) {
    return (
      <div>
        <h4 className='text-center'>
          Bem vindo! Esperamos que você se divirta e faça novas amizades. Estamos no{' '}
          <span
            onClick={() => openDiscord()}
            style={{
              color: '#7289DA',
            }}
          >
            Discord
          </span>{' '}
          para ajudar no que for possível, então qualquer dúvida ou sugestão entre em contato.
          <br />~ Grato Arca
        </h4>
        <div className='flex items-center justify-center flex-col'>
          <Button onClick={handleStartGame}>Entrar</Button>
          <Button onClick={handleChangeAccount}>Trocar de conta</Button>
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
        <h3 style={{}}>Selecione seu personagem</h3>
        <Button onClick={handleStartGame} disabled={!selectedCharId || newCharInLoading}>
          Entrar
        </Button>
        <Button
          onClick={handleCreateNewChar}
          loadingMessage='Criando novo personagem...'
          isLoading={newCharInLoading}
        >
          Novo personagem
        </Button>
        <Button onClick={handleChangeAccount} disabled={newCharInLoading}>
          Trocar de conta
        </Button>
        <AutoLoginCheckbox />
      </div>
    </div>
  );
}

function AutoLoginCheckbox() {
  const [checked, setChecked] = useState(() => localStorage.getItem('autoLogin') === 'true');

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
      <label htmlFor='auto-login'>Continuar conectado</label>
    </div>
  );
}
