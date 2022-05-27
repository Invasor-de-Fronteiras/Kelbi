import React, { useMemo } from 'react';

import Palamute from '../assets/loading/palamute.gif';
import Palico from '../assets/loading/palico.gif';
import { openBrowser } from '../utils/launcher';
import { randomArr } from '../utils/util';

export function Loading() {
  const data = useMemo(
    () =>
      randomArr([
        {
          img: Palamute,
          width: '100',
          height: '100',
          alt: 'Palamute',
          href: 'https://www.youtube.com/watch?v=0Uk-gytGg94',
        },
        {
          img: Palico,
          width: '100%',
          height: '100%',
        },
      ]),
    [],
  );

  return (
    <div
      className='hover-help'
      style={{
        position: 'absolute',
        bottom: 0,
        right: 0,
      }}
      onClick={() => (data.href ? openBrowser(data.href) : null)}
    >
      <img id='img-bg' src={data.img} height={data.height} width={data.width} />
    </div>
  );
}
