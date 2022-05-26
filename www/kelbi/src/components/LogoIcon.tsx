import React from 'react';
import LogoImg from '../assets/logo.png';

export function LogoIcon() {
  const onHold = (active: boolean) => {
    // TODO: move to launcher.ts
    //@ts-ignore
    window.external.beginDrag(active);
  };

  return (
    <img
      src={LogoImg}
      width='303'
      height='114'
      onMouseDown={() => onHold(true)}
      onMouseUp={() => onHold(false)}
      onMouseLeave={() => onHold(false)}
      onTouchStart={() => onHold(true)}
      onTouchEnd={() => onHold(false)}
    />
  );
}
