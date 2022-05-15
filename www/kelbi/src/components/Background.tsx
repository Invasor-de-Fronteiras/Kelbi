import React, { useMemo } from 'react';

import BgImg1 from '../assets/background/1.jpg';
import BgImg2 from '../assets/background/2.jpg';
import BgImg3 from '../assets/background/3.jpg';
import BgImg4 from '../assets/background/4.jpg';
import BgImg5 from '../assets/background/5.jpg';
import BgImg6 from '../assets/background/6.jpg';
import BgImg7 from '../assets/background/7.jpg';
import BgImg8 from '../assets/background/8.jpg';
import BgImg9 from '../assets/background/9.jpg';
import BgImg10 from '../assets/background/10.jpg';
import BgImg11 from '../assets/background/11.jpg';
import BgImg12 from '../assets/background/12.jpg';
import BgImg13 from '../assets/background/13.jpg';
import BgImg14 from '../assets/background/14.jpg';
import BgImg15 from '../assets/background/15.jpg';
import BgImg16 from '../assets/background/16.jpg';
import BgImg17 from '../assets/background/17.jpg';
import BgImg18 from '../assets/background/18.jpg';
import BgImg19 from '../assets/background/19.jpg';

function randomArr<T>(arr: T[]) {
  return arr[Math.floor(Math.random() * arr.length)];
}

const images = [
  BgImg1,
  BgImg2,
  BgImg3,
  BgImg4,
  BgImg5,
  BgImg6,
  BgImg7,
  BgImg8,
  BgImg9,
  BgImg10,
  BgImg11,
  BgImg12,
  BgImg13,
  BgImg14,
  BgImg15,
  BgImg16,
  BgImg17,
  BgImg18,
  BgImg19,
];

export function Background() {
  const bgImg = useMemo(() => randomArr(images), []);

  return (
    <div>
      <img
        id='img-bg'
        src={bgImg}
        height='600px'
        // width='100%'
        width='1124px'
        style={{
          position: 'absolute',
          top: 0,
          left: 0,
          zIndex: -1,
        }}
      />
    </div>
  );
}
