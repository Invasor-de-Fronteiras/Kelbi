import React from 'react';
import { Character } from '../utils/launcher';
import { normalizeDate } from '../utils/util';

import './CharacterCard.css';

interface CharacterCardProps {
  char: Character;
  tabIndex: number;
  onSelect: () => void;
  isSelected: boolean;
}

export function CharacterCard({ char, tabIndex, isSelected, onSelect }: CharacterCardProps) {
  return (
    <div
      className='character-card'
      role='option'
      tabIndex={tabIndex}
      onClick={onSelect}
      style={{
        backgroundColor: isSelected ? '#f83' : '#fff',
      }}
    >
      <div className='flex flex-row justify-between items-center'>
        <span className='char-weapon'>{char.weapon}</span>
      </div>
      <h3 className='char-name'>{char.name}</h3>
      <div
        style={{
          marginTop: '5px',
          maxWidth: '200px',
        }}
      >
        <div className='flex flex-row justify-between items-center'>
          <span className='flex-1'>HR {char.HR}</span>
          {char.GR ? <span className='flex-1'>GR {char.GR}</span> : null}
        </div>
        <div className='flex flex-row justify-between items-center'>
          <span className='flex-1'>Last Login</span>
          <span className='flex-1'>{normalizeDate(char.lastLogin)}</span>
        </div>
      </div>
    </div>
  );
}
