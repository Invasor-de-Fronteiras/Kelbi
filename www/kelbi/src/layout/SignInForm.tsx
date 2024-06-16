import React, { useEffect } from 'react';
import { FieldCheckbox } from '../components/FieldCheckbox';
import { FieldInput } from '../components/FieldInput';
import { Button } from '../components/Button';
import { FormikProvider, useFormik } from 'formik';
import * as Yup from 'yup';
import { RiErrorWarningFill } from 'react-icons/ri';

import './SignInForm.css';
import { useLogin } from '../hooks/useLogin';
import { useTranslate } from '../i18n/useTranslate';
import { LocaleKeys } from '../i18n/TranslateProvider';

const validationSchema = Yup.object().shape({
  username: Yup.string().required('login_field_required'),
  password: Yup.string().required('login_field_required'),
  autoLogin: Yup.boolean().required('login_field_required'),
});

export function SignInForm() {
  const { t } = useTranslate();

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
    validateOnMount: true,
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
        <h1>{t('login_title')}</h1>
        <FieldInput
          placeholder={t('login_username')}
          type='text'
          name='username'
          isRequired
          disabled={isLoading}
        />
        <FieldInput
          placeholder={t('login_password')}
          type='password'
          isRequired
          name='password'
          disabled={isLoading}
        />
        <FieldCheckbox name='autoLogin' disabled={isLoading}>
          {t('login_keep')}
        </FieldCheckbox>
        <Button
          type='submit'
          isLoading={isLoading}
          disabled={formik.isValidating || !formik.isValid}
        >
          {t('login_btn')}
        </Button>
        {error && <SignInError error={error} />}
      </form>
    </FormikProvider>
  );
}

function SignInError({ error }: { error: Error }) {
  const { t } = useTranslate();

  return (
    <div
      style={{
        marginTop: '10px',
        backgroundColor: '#B00020',
        padding: '10px',
      }}
    >
      <div
        style={{
          display: 'flex',
          flexDirection: 'row',
          alignItems: 'center',
          justifyContent: 'space-between',
        }}
      >
        <span
          style={{
            color: '#fff',
          }}
        >
          {t('login_error_title')}
        </span>
        <RiErrorWarningFill size={20} color='#fff' />
      </div>
      <span
        style={{
          color: '#fff',
          marginTop: '5px',
          fontSize: '14px',
        }}
      >
        {t(error)}
      </span>
    </div>
  );
}
