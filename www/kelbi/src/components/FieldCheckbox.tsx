import React from 'react';

interface FieldCheckboxProps {
  name: string;
  label?: string;
  value?: boolean;
  onChange?: (newValue: boolean) => void;
}

export function FieldCheckbox({
  name,
  value,
  onChange = () => null,
  ...props
}: FieldCheckboxProps) {
  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'row',
        alignItems: 'center',
      }}>
      <input
        type='checkbox'
        id={name}
        checked={value}
        onChange={(e) => onChange(e.target.checked)}
        {...props}
      />
      <label
        htmlFor={name}
        style={{
          marginLeft: '10px',
        }}>
        manter login
      </label>
    </div>
  );
}
