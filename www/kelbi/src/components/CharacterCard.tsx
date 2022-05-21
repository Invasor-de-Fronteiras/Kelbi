import React from 'react';
import { Character } from '../utils/launcher';
import { normalizeDate } from '../utils/util';

interface CharacterCardProps {
  char: Character;
}

export function CharacterCard({ char }: CharacterCardProps) {
  return (
    <div
      style={{
        marginBottom: '10px',
        padding: '2px',
      }}>
      <span
        style={{
          textTransform: 'uppercase',
          fontSize: '.9rem',
        }}>
        {char.weapon}
      </span>
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
