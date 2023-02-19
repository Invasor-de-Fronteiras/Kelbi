import React from 'react';
import { useTranslate } from '../i18n/useTranslate';
import { playClickSong } from '../utils/songs';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children: React.ReactNode;
  isLoading?: boolean;
  loadingMessage?: string;
}

export function Button({
  isLoading,
  children,
  disabled,
  onClick,
  loadingMessage,
  ...props
}: ButtonProps) {
  const { t } = useTranslate();

  return (
    <button
      {...props}
      onClick={(e) => {
        playClickSong();
        onClick?.(e);
      }}
      disabled={disabled || isLoading}
    >
      {isLoading ? <span className='loading'>{loadingMessage ?? t('loading')}</span> : children}
    </button>
  );
}
