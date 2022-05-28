import classNames from 'classnames';
import { useField } from 'formik';
import React from 'react';

interface FieldCheckboxProps {
  name: string;
  children: React.ReactNode;
  disabled?: boolean;
}

export function FieldCheckbox({ name, children, disabled }: FieldCheckboxProps) {
  const [field] = useField({ name, type: 'checkbox' });

  return (
    <div className={classNames('checkbox-group', { disabled })}>
      <input type='checkbox' id={name} {...field} />
      <label htmlFor={name}>{children}</label>
    </div>
  );
}
