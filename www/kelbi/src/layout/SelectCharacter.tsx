import React from 'react';

import BabaroWalking from '../assets/loading/banbaro-walking.gif';
import { CharacterCard } from '../components/CharacterCard';
import { useGetCharacters } from '../hooks/useGetCharacters';
export function SelectCharacter() {
  const { characters, isNewAccount, loading } = useGetCharacters();

  if (loading) return <SelectCharacterLoading />;
  return (
    <div>
      <div
        style={{
          overflowY: 'auto',
          maxHeight: '230px',
          marginBottom: '30px',
        }}>
        {characters.map((char) => (
          <CharacterCard char={char} key={char.uid} />
        ))}
      </div>
      <div
        style={{
          display: 'flex',
          flexDirection: 'row',
          justifyContent: 'space-between',
        }}>
        <button style={{ flex: 1, marginRight: 5 }}>Novo Char</button>
        <button style={{ flex: 1, marginLeft: 5 }}>Sair da conta</button>
      </div>
      <div
        className='checkbox-group'
        style={{
          marginTop: '10px',
        }}>
        <input type='checkbox' id='auto-login' />
        <label htmlFor='auto-login'>Manter login</label>
      </div>
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
      <img id='img-bg' src={BabaroWalking} height='100%' width='100%' />
      <span
        style={{
          marginTop: '5px',
        }}>
        Buscando seus dados...
      </span>
    </div>
  );
}
