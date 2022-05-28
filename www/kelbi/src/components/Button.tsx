import React from 'react';
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
  loadingMessage = 'Carregando...',
  ...props
}: ButtonProps) {
  return (
    <button
      {...props}
      onClick={(e) => {
        playClickSong();
        onClick?.(e);
      }}
      disabled={disabled || isLoading}
    >
      {isLoading ? <span className='loading'>{loadingMessage}</span> : children}
    </button>
  );
}
