import classNames from 'classnames';
import { ErrorMessage, useField } from 'formik';
import React, { ComponentProps } from 'react';
import { playDigitSong } from '../utils/songs';

interface FieldInputProps extends ComponentProps<'input'> {
  name: string;
  isRequired?: boolean;
  label?: string;
}

export function FieldInput({ name, isRequired, label, ...props }: FieldInputProps) {
  const [field, meta] = useField(name);

  const hasAnErrorAndHasBeenTouched = !!meta.error && meta.touched;

  return (
    <div className='input-group'>
      <label htmlFor={name}>{label}</label>
      <input
        id={name}
        className={classNames('input', {
          'input-error': hasAnErrorAndHasBeenTouched,
        })}
        {...props}
        required={isRequired}
        {...field}
        onChange={(e) => {
          field.onChange(e);
          playDigitSong();
        }}
      />
      <ErrorMessage name={name} render={(err) => <span className='error-message'>{err}</span>} />
    </div>
  );
}
