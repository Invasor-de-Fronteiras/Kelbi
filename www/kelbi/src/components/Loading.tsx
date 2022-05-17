import React, { useMemo } from 'react';

import Palamute from '../assets/loading/palamute.gif';
import Palico from '../assets/loading/palico.gif';
import { randomArr } from '../utils/util';

export function Loading() {
  const data = useMemo(
    () =>
      randomArr([
        { img: Palamute, width: '100', height: '100' },
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
      style={{
        position: 'absolute',
        bottom: 0,
        right: 0,
      }}>
      <img id='img-bg' src={data.img} height={data.height} width={data.width} />
    </div>
  );
}
