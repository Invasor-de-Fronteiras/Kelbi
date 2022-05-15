import React from 'react';

export function FieldCheckbox() {
  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'row',
        alignItems: 'center',
      }}>
      <label
        style={{
          marginRight: '10px',
        }}>
        Remember me
      </label>
      <input type='checkbox' />
    </div>
  );
}
