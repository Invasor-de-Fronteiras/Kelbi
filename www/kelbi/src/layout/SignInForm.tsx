import React from 'react';
import { useLauncher } from '../context/LauncherContext';
import { FieldCheckbox } from '../components/FieldCheckbox';
import { FieldInput } from '../components/FieldInput';
import { Button } from '../components/Button';
import { FormikProvider, useFormik } from 'formik';
import * as Yup from 'yup';

import './SignInForm.css';

interface FormValues {
  accountId: string;
  password: string;
  rememberMe: boolean;
}

const validationSchema = Yup.object().shape({
  email: Yup.string().required('Campo obrigatório.'),
  password: Yup.string().min(6, 'Senha muito curta.').required('Campo obrigatório.'),
  rememberMe: Yup.boolean().required('Campo obrigatório.'),
});

const startButtonInDisabledState = {
  initialErrors: {
    rememberMe: 'ignore',
  },
};

export function SignInForm() {
  const { isLoading, setIsLoading } = useLauncher();

  const formik = useFormik({
    initialValues: {
      email: '',
      password: '',
      rememberMe: false,
    },
    validationSchema,
    onSubmit: async (data) => {
      console.log(data);
      setIsLoading(true);
      await new Promise((resolve) => setTimeout(resolve, 1000));
      setIsLoading(false);
    },
    ...startButtonInDisabledState,
  });

  return (
    <FormikProvider value={formik}>
      <form id='sign-in' onSubmit={formik.handleSubmit}>
        <h1>Fazer login </h1>
        <FieldInput
          placeholder='Nome de usuário'
          type='text'
          name='email'
          isRequired
          disabled={isLoading}
        />
        <FieldInput
          placeholder='Senha'
          type='password'
          isRequired
          name='password'
          disabled={isLoading}
        />
        <FieldCheckbox name='rememberMe' disabled={isLoading}>
          manter login
        </FieldCheckbox>
        <Button
          type='submit'
          isLoading={isLoading}
          disabled={formik.isValidating || !formik.isValid}>
          Entrar
        </Button>
      </form>
    </FormikProvider>
  );
}
