import React, { ComponentProps } from 'react';
// import { useField } from 'formik';

interface FieldInputProps extends ComponentProps<'input'> {
  name: string;
  isRequired?: boolean;
  label?: string;
}

export function FieldInput({ name, isRequired, label, ...props }: FieldInputProps) {
  // const [field] = useField(name);

  // const hasAnErrorAndHasBeenTouched = !!meta.error && meta.touched;

  return (
    <div className='input-group'>
      <label htmlFor={name}>{label}</label>
      <input id={name} className='input' {...props} required={isRequired} />
    </div>
  );
}
