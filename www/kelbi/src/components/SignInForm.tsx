import React from 'react';
import { FieldCheckbox } from './FieldCheckbox';
import { FieldInput } from './FieldInput';
import { PasswordInput } from './PasswordInput';

import './SignInForm.css';

export function SignInForm() {
  return (
    <form id='sign-in'>
      <h1>Sign In</h1>
      <FieldInput />
      <PasswordInput />
      <FieldCheckbox />
      <button type='submit'>Sign In</button>
    </form>
  );
}
