import React from 'react';
import LogoImg from '../assets/logo.png';

export function LogoIcon() {
  return (
    <img
      src={LogoImg}
      width='303'
      height='114'
      onMouseDown={() => window.external.beginDrag(true)}
      onMouseUp={() => window.external.beginDrag(false)}
      onMouseLeave={() => window.external.beginDrag(false)}
      onTouchStart={() => window.external.beginDrag(true)}
      onTouchEnd={() => window.external.beginDrag(false)}
    />
  );
}
