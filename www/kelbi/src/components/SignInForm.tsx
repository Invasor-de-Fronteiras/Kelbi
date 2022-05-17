import React from 'react';
import { useLauncher } from '../context/LauncherContext';
import { FieldCheckbox } from './FieldCheckbox';
import { FieldInput } from './FieldInput';
import './SignInForm.css';
import { Button } from './Button';
import { FormikProvider, useFormik } from 'formik';

interface FormValues {
  accountId: string;
  password: string;
  rememberMe: boolean;
}

export function SignInForm() {
  const { isLoading, setIsLoading } = useLauncher();

  const formik = useFormik({
    initialValues: {
      accountId: '',
      password: '',
      rememberMe: false,
    },
    onSubmit: async () => {
      setIsLoading(true);
      await new Promise((resolve) => setTimeout(resolve, 1000));
      setIsLoading(false);
    },
  });

  return (
    <FormikProvider value={formik}>
      <form id='sign-in' onSubmit={formik.handleSubmit}>
        <h1>Fazer login </h1>
        <FieldInput
          placeholder='Nome de usuário'
          type='text'
          isRequired
          name='accountId'
          disabled={isLoading}
        />
        <FieldInput
          placeholder='Senha'
          type='password'
          isRequired
          name='password'
          disabled={isLoading}
        />
        <FieldCheckbox name='rememberMe' disabled={isLoading} />
        <Button type='submit' isLoading={isLoading}>
          Entrar
        </Button>
      </form>
    </FormikProvider>
  );
}
