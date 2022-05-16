import React, { useState } from 'react';
import { useLauncher } from '../context/LauncherContext';
import { FieldCheckbox } from './FieldCheckbox';
import { FieldInput } from './FieldInput';
import './SignInForm.css';
import { Button } from './Button';

interface FormValues {
  accountId: string;
  password: string;
  rememberMe: boolean;
}

export function SignInForm() {
  const { isLoading, setIsLoading } = useLauncher();

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
    setIsLoading(true);
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
        disabled={isLoading}
      />
      <FieldInput
        placeholder='Senha'
        type='password'
        isRequired
        {...makeInputProps('password')}
        disabled={isLoading}
      />
      <FieldCheckbox {...makeCheckboxProps('rememberMe')} disabled={isLoading} />
      <Button isLoading={isLoading}>Entrar</Button>
    </form>
  );
}
