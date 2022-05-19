import React, { useEffect } from 'react';
import { FieldCheckbox } from '../components/FieldCheckbox';
import { FieldInput } from '../components/FieldInput';
import { Button } from '../components/Button';
import { FormikProvider, useFormik } from 'formik';
import * as Yup from 'yup';
import { RiErrorWarningFill } from 'react-icons/ri';

import './SignInForm.css';
import { useLogin } from '../hooks/useLogin';

interface FormValues {
  username: string;
  password: string;
  autoLogin: boolean;
}

const validationSchema = Yup.object().shape({
  username: Yup.string().required('Campo obrigatório.'),
  password: Yup.string().min(6, 'Senha muito curta.').required('Campo obrigatório.'),
  autoLogin: Yup.boolean().required('Campo obrigatório.'),
});

export function SignInForm() {
  const { error, isLoading, mutate } = useLogin({
    onSuccess: (input) => {
      localStorage.setItem('username', input.username);
      localStorage.setItem('password', input.password);
      localStorage.setItem('autoLogin', String(input.autoLogin));
    },
  });

  const initialValues = {
    username: localStorage.getItem('username') ?? '',
    password: localStorage.getItem('password') ?? '',
    autoLogin: localStorage.getItem('autoLogin') === 'true',
  };

  const formik = useFormik({
    initialValues,
    validationSchema,
    isInitialValid: validationSchema.isValidSync(initialValues),
    onSubmit: (data) => mutate(data),
  });

  useEffect(() => {
    if (validationSchema.isValidSync(initialValues) && initialValues.autoLogin) {
      mutate(initialValues);
    }
  }, []);

  return (
    <FormikProvider value={formik}>
      <form id='sign-in' onSubmit={formik.handleSubmit}>
        <h1>Fazer login</h1>
        <FieldInput
          placeholder='Nome de usuário'
          type='text'
          name='username'
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
        <FieldCheckbox name='autoLogin' disabled={isLoading}>
          manter login
        </FieldCheckbox>
        <Button
          type='submit'
          isLoading={isLoading}
          disabled={formik.isValidating || !formik.isValid}>
          Entrar
        </Button>
        {error && <SignInError error={error} />}
      </form>
    </FormikProvider>
  );
}

function SignInError({ error }: { error: Error }) {
  return (
    <div
      style={{
        marginTop: '10px',
        backgroundColor: '#B00020',
        padding: '10px',
      }}>
      <div
        style={{
          display: 'flex',
          flexDirection: 'row',
          alignItems: 'center',
          justifyContent: 'space-between',
        }}>
        <span
          style={{
            color: '#fff',
          }}>
          Oops!
        </span>
        <RiErrorWarningFill size={20} color='#fff' />
      </div>
      <span
        style={{
          color: '#fff',
          marginTop: '5px',
          fontSize: '14px',
        }}>
        {error?.message}
      </span>
    </div>
  );
}
