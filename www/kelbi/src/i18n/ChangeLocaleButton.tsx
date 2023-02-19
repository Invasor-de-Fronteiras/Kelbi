import React, { useMemo, useState } from 'react';
import { IoIosArrowUp } from 'react-icons/io';
import { BsCheck } from 'react-icons/bs';

import './ChangeLocaleButton.css';
import { useTranslate } from './useTranslate';

import enUS from '../assets/flags/en-US.png';
import esES from '../assets/flags/es-ES.png';
import ptBR from '../assets/flags/pt-BR.png';
import ja from '../assets/flags/ja.png';

const languages = [
  {
    id: 'en-US',
    title: 'English, USA',
    icon: enUS,
  },
  {
    id: 'es-ES',
    title: 'Español',
    icon: esES,
  },
  {
    id: 'pt-BR',
    title: 'Português do Brasil',
    icon: ptBR,
  },
  {
    id: 'ja',
    title: '日本語',
    icon: ja,
  },
] as const;

export const ChangeLocaleButton = () => {
  const { locale, setLocale } = useTranslate();
  const [showModal, setShowModal] = useState(true);

  const lang = useMemo(() => languages.find((lang) => lang.id === locale), [locale]);

  return (
    <div id='locale-content'>
      <div
        id='open-locale-btn'
        onClick={() => setShowModal(true)}
        style={{
          padding: '3px',
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          marginLeft: '10px',
        }}
      >
        <div
          style={{
            maxWidth: '25px',
            maxHeight: '25px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            marginRight: '4px',
          }}
        >
          <img src={lang.icon} width='100%' height='100%' />
        </div>
        <IoIosArrowUp size={15} />
      </div>
      {showModal && (
        <div
          style={{
            position: 'absolute',
            top: 0,
            bottom: 0,
            left: 0,
            right: 0,
            width: '100vw',
            height: '100vh',
            backgroundColor: '#000',
            opacity: 0.5,
          }}
          onClick={() => setShowModal(false)}
        />
      )}
      {showModal && (
        <div id='locale-modal'>
          <ul>
            {languages.map((lang) => {
              const selected = lang.id === locale;

              return (
                <li
                  key={lang.id}
                  className={selected ? 'selected' : ''}
                  onClick={
                    selected
                      ? undefined
                      : () => {
                          setLocale(lang.id);
                          setShowModal(false);
                        }
                  }
                >
                  <div
                    style={{
                      maxWidth: '25px',
                      maxHeight: '25px',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      marginRight: '4px',
                    }}
                  >
                    <img src={lang.icon} width='100%' height='100%' />
                  </div>
                  {lang.title} {selected && <BsCheck color='#fff' />}
                </li>
              );
            })}
          </ul>
        </div>
      )}
    </div>
  );
};
