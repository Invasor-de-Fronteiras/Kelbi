import React, { useState } from 'react';

import BanbaroWalking from '../assets/loading/banbaro-walking.gif';
import { Button } from '../components/Button';
import { CharacterCard } from '../components/CharacterCard';
import { useGetCharacters } from '../hooks/useGetCharacters';
export function SelectCharacter() {
  const { characters, isNewAccount, loading, username } = useGetCharacters();

  if (loading) return <SelectCharacterLoading />;

  if (isNewAccount) {
    return (
      <div>
        <h4 className='text-center'>
          Olá {username}, Queremos lhe desejar boa sorte e boas-vindas, que você não desista do jogo
          antes do G Rank!
        </h4>
        <div className='flex items-center justify-center flex-col'>
          <Button>Entrar</Button>
          <Button>Trocar de conta</Button>
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
      <div
        style={{
          display: 'flex',
          flexDirection: 'row',
          justifyContent: 'space-between',
        }}>
        <Button>Novo personagem</Button>
        <Button>Trocar de conta</Button>
        <AutoLoginCheckbox />
      </div>
    </div>
  );
}

function AutoLoginCheckbox() {
  const [checked, setChecked] = useState(() => localStorage.getItem('autoLogin') === 'true');

  function handleChange(e) {
    localStorage.setItem('autoLogin', e.target.checked);
    setChecked(e.target.checked);
  }

  return (
    <div
      className='checkbox-group'
      style={{
        marginTop: '10px',
      }}>
      <input type='checkbox' id='auto-login' onChange={handleChange} checked={checked} />
      <label htmlFor='auto-login'>Continuar conectado</label>
    </div>
  );
}
function SelectCharacterLoading() {
  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
      }}>
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
