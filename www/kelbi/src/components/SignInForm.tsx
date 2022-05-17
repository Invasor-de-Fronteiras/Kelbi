import React from 'react';
import { useLauncher } from '../context/LauncherContext';
import { FieldCheckbox } from './FieldCheckbox';
import { FieldInput } from './FieldInput';
import './SignInForm.css';
import { Button } from './Button';
import { FormikProvider, useFormik } from 'formik';
import * as Yup from 'yup';

interface FormValues {
  accountId: string;
  password: string;
  rememberMe: boolean;
}

const validationSchema = Yup.object().shape({
  email: Yup.string().email('E-mail inválido.').required('Campo obrigatório.'),
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
    onSubmit: async () => {
      setIsLoading(true);
      await new Promise((resolve) => setTimeout(resolve, 1000));
      setIsLoading(false);
    },
    ...startButtonInDisabledState,
  });

  console.log(formik);

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
        <FieldCheckbox name='rememberMe' disabled={isLoading} />
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
