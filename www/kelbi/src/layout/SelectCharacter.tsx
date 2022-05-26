import React, { useState } from 'react';

import BarbaryWalking from '../assets/loading/banbaro-walking.gif';
import { Button } from '../components/Button';
import { CharacterCard } from '../components/CharacterCard';
import { useCreateCharacter } from '../hooks/useCreateCharacter';
import { useGetCharacters } from '../hooks/useGetCharacters';
import { startGame } from '../utils/launcher';

// TODO: move new accounts to an own file
export function SelectCharacter() {
  const { mutate: handleCreateNewChar } = useCreateCharacter();

  const { characters, isNewAccount, newAccountUID, loading } = useGetCharacters();
  const [selectedCharId, setSelectedCharId] = useState('');

  const handleStartGame = () => {
    // START GAME FOR NEW ACCOUNT
    startGame(selectedCharId ?? newAccountUID);

    // START GAME FOR SELECTED CHARACTER
  };

  const handleChangeAccount = () => {
    localStorage.removeItem('autoLogin');
    localStorage.removeItem('password');
    localStorage.removeItem('username');
    window.external.restartMhf();
  };

  if (loading) return <SelectCharacterLoading />;

  if (isNewAccount) {
    return (
      <div>
        <h4 className='text-center'>
          Olá {window.external.getUserId()}, Queremos lhe desejar boa sorte e boas-vindas, que você
          não desista do jogo antes do G Rank!
        </h4>
        <div className='flex items-center justify-center flex-col'>
          <Button onClick={handleStartGame}>Entrar</Button>
          <Button onClick={handleChangeAccount}>Trocar de conta</Button>
        </div>
      </div>
    );
  }

  return (
    <div>
      <div
        style={{
          overflowY: 'auto',
          maxHeight: '230px',
          marginBottom: '30px',
        }}>
        {characters.map((char, index) => (
          <CharacterCard
            char={char}
            key={char.uid}
            tabIndex={index}
            isSelected={char.uid === selectedCharId}
            onSelect={() => setSelectedCharId(char.uid)}
          />
        ))}
      </div>
      <div className='flex flex-col justify-between'>
        <Button onClick={handleStartGame}>Entrar</Button>
        <Button onClick={handleCreateNewChar}>Novo personagem</Button>
        <Button onClick={handleChangeAccount}>Trocar de conta</Button>
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
      onClick={() => handleChange(!checked)}>
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

function SelectCharacterLoading() {
  return (
    <div className='flex items-center justify-center flex-col'>
      <img id='img-bg' src={BarbaryWalking} height='100%' width='100%' />
      <span
        style={{
          marginTop: '5px',
        }}>
        Buscando seus dados...
      </span>
    </div>
  );
}
