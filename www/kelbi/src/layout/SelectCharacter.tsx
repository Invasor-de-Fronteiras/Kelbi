import React, { useState } from 'react';

import BanbaroWalking from '../assets/loading/banbaro-walking.gif';
import { Button } from '../components/Button';
import { CharacterCard } from '../components/CharacterCard';
import { useGetCharacters } from '../hooks/useGetCharacters';
import { getCharacters, isNeAccountChar } from '../utils/launcher';
import { findByEnd } from '../utils/util';

// TODO: move new accounts to an own file
export function SelectCharacter() {
  const { characters, isNewAccount, newAccountUID, loading, username } = useGetCharacters();

  const handleStartGame = () => {
    // START GAME FOR NEW ACCOUNT
    // @ts-ignore
    window.external.selectCharacter(newAccountUID, newAccountUID);
    // @ts-ignore
    window.external.exitLauncher();

    // START GAME FOR SELECTED CHARACTER
  };

  const handleChangeAccount = () => {
    localStorage.removeItem('autoLogin');
    localStorage.removeItem('password');
    localStorage.removeItem('username');
    // @ts-ignore
    window.external.restartMhf();
  };

  const handleCreateNewChar = () => {
    const username = localStorage.getItem('username');
    const password = localStorage.getItem('password');

    // check if has a new char for not need create other

    // @ts-ignore
    window.external.loginCog(username + '+', password, password);

    setInterval(() => {
      //Todo: check loading by lastAuthResult
      const chars = getCharacters();
      const newChar = findByEnd(chars, isNeAccountChar);

      if (!newChar) {
        // Todo: throw error or continue in loading
        return;
      }

      // @ts-ignore
      window.external.selectCharacter(newAccountUID, newAccountUID);
      // @ts-ignore
      window.external.exitLauncher();
    }, 100);
  };

  const handleDeleteChar = () => {
    alert(`${username} is ready to delete character!`);
  };

  if (loading) return <SelectCharacterLoading />;

  if (isNewAccount) {
    return (
      <div>
        <h4 className='text-center'>
          Olá {username}, Queremos lhe desejar boa sorte e boas-vindas, que você não desista do jogo
          antes do G Rank!
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
        {/* <CharacterCard
          tabIndex={1}
          char={{
            GR: 33,
            name: 'Banbaro',
            HR: 0,
            gender: 'Famele',
            lastLogin: new Date(),
            uid: 'wra',
            weapon: 'Bow',
          }}
        />
        <CharacterCard
          tabIndex={2}
          char={{
            GR: 33,
            name: 'Banbaro',
            HR: 0,
            gender: 'Famele',
            lastLogin: new Date(),
            uid: 'wra',
            weapon: 'Bow',
          }}
        />
        <CharacterCard
          tabIndex={3}
          char={{
            GR: 33,
            name: 'Banbaro',
            HR: 0,
            gender: 'Famele',
            lastLogin: new Date(),
            uid: 'wra',
            weapon: 'Bow',
          }}
        /> */}
        {characters.map((char, index) => (
          <CharacterCard char={char} key={char.uid} tabIndex={index} />
        ))}
      </div>
      <div className='flex flex-row justify-between'>
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
      <img id='img-bg' src={BanbaroWalking} height='100%' width='100%' />
      <span
        style={{
          marginTop: '5px',
        }}>
        Buscando seus dados...
      </span>
    </div>
  );
}
