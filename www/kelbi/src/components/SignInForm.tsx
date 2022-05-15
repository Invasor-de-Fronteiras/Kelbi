import React, { useState } from 'react';
import { FieldCheckbox } from './FieldCheckbox';
import { FieldInput } from './FieldInput';

import './SignInForm.css';

interface FormValues {
  accountId: string;
  password: string;
  rememberMe: boolean;
}

export function SignInForm() {
  const [form, setForm] = useState({
    values: {
      accountId: '',
      password: '',
      rememberMe: false,
    },
    errors: {
      accountId: '',
      password: '',
    },
  });

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
  };

  const makeInputProps = (name: Exclude<keyof FormValues, 'rememberMe'>) => {
    const { values } = form;

    return {
      name,
      value: values[name],
      onChange: (e: React.ChangeEvent<HTMLInputElement>) => {
        const { value } = e.target;
        setForm((prev) => ({
          errors: { ...prev.errors, [name]: null },
          values: { ...prev.values, [name]: value },
        }));
      },
    };
  };

  const makeCheckboxProps = (name: 'rememberMe') => {
    const { values } = form;

    return {
      name,
      value: values[name],
      onChange: (v: boolean) => {
        setForm((prev) => ({
          errors: { ...prev.errors, [name]: null },
          values: {
            ...prev.values,
            [name]: v,
          },
        }));
      },
    };
  };

  return (
    <form id='sign-in' onSubmit={onSubmit}>
      <h1>Fazer login </h1>
      <FieldInput
        placeholder='Nome de usuário'
        type='text'
        isRequired
        {...makeInputProps('accountId')}
      />
      <FieldInput placeholder='Senha' type='password' isRequired {...makeInputProps('password')} />
      <FieldCheckbox {...makeCheckboxProps('rememberMe')} />
      <button type='submit'>Entrar</button>
    </form>
  );
}
