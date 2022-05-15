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
      <label
        htmlFor={name}
        style={{
          marginRight: '10px',
        }}>
        Remember me
      </label>
      <input
        type='checkbox'
        id={name}
        checked={value}
        onChange={(e) => onChange(e.target.checked)}
        {...props}
      />
    </div>
  );
}
