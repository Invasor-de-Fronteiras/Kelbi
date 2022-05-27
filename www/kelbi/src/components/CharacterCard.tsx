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
        marginBottom: '10px',
        width: '200px',
        padding: '10px',
      }}
    >
      <div
        style={{
          display: 'flex',
          flexDirection: 'row',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        <span
          style={{
            textTransform: 'uppercase',
            fontSize: '.9rem',
          }}
        >
          {char.weapon}
        </span>
      </div>
      <h3
        style={{
          fontSize: '1.5rem',
          fontWeight: 'bold',
          textTransform: 'uppercase',
        }}
      >
        {char.name}
      </h3>
      <div
        style={{
          marginTop: '5px',
          maxWidth: '200px',
        }}
      >
        <div
          style={{
            display: 'flex',
            flexDirection: 'row',
            justifyContent: 'space-between',
          }}
        >
          <span style={{ flex: 1 }}>HR {char.HR}</span>
          {char.GR ? <span style={{ flex: 1 }}>GR {char.GR}</span> : null}
        </div>
        <div
          style={{
            display: 'flex',
            flexDirection: 'row',
            justifyContent: 'space-between',
          }}
        >
          <span style={{ flex: 1 }}>Last Login</span>
          <span style={{ flex: 1 }}>{normalizeDate(char.lastLogin)}</span>
        </div>
      </div>
    </div>
  );
}
