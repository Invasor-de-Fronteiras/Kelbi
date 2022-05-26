import React from 'react';
import { Character } from '../utils/launcher';
import { normalizeDate } from '../utils/util';
import { RiDeleteBinLine } from 'react-icons/ri';

interface CharacterCardProps {
  char: Character;
  tabIndex: number;
  onSelect: () => void;
  isSelected: boolean;
}

export function CharacterCard({ char, tabIndex, isSelected, onSelect }: CharacterCardProps) {
  return (
    <div
      role='option'
      tabIndex={tabIndex}
      onClick={onSelect}
      style={{
        backgroundColor: isSelected ? '#f83' : '#fff',
        marginBottom: '10px',
        padding: '3px',
      }}>
      <div
        style={{
          display: 'flex',
          flexDirection: 'row',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}>
        <span
          style={{
            textTransform: 'uppercase',
            fontSize: '.9rem',
          }}>
          {char.weapon}
        </span>
        <RiDeleteBinLine color='#B00020' className='hover-pointer' />
      </div>
      <h3
        style={{
          fontSize: '1.5rem',
          fontWeight: 'bold',
          textTransform: 'uppercase',
        }}>
        {char.name}
      </h3>
      <div
        style={{
          marginTop: '5px',
          maxWidth: '200px',
        }}>
        <div
          style={{
            display: 'flex',
            flexDirection: 'row',
            justifyContent: 'space-between',
          }}>
          <span style={{ flex: 1 }}>HR {char.HR}</span>
          {char.GR && <span style={{ flex: 1 }}>GR {char.GR}</span>}
        </div>
        <div
          style={{
            display: 'flex',
            flexDirection: 'row',
            justifyContent: 'space-between',
          }}>
          <span style={{ flex: 1 }}>Last Login</span>
          <span style={{ flex: 1 }}>{normalizeDate(char.lastLogin)}</span>
        </div>
      </div>
    </div>
  );
}
