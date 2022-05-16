import React, { useMemo } from 'react';

import Palamute from '../assets/loading/palamute.gif';
import Palico from '../assets/loading/palico.gif';
import { randomArr } from '../utils/util';

export function Loading() {
  const loadingImg = useMemo(() => randomArr([Palamute, Palico]), []);

  return (
    <div
      style={{
        position: 'absolute',
        bottom: 0,
        right: 0,
      }}>
      <img id='img-bg' src={loadingImg} height='100' width='100' />
      {/* {message && <span className='loading-message'>{message}</span>} */}
    </div>
  );
}
