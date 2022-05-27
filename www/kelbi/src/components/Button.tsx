import React from 'react';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children: React.ReactNode;
  isLoading?: boolean;
  loadingMessage?: string;
}

export function Button({
  isLoading,
  children,
  disabled,
  loadingMessage = 'Carregando...',
  ...props
}: ButtonProps) {
  return (
    <button {...props} disabled={disabled || isLoading}>
      {isLoading ? <span className='loading'>{loadingMessage}</span> : children}
    </button>
  );
}
