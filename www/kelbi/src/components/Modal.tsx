import React from 'react';
import { Button } from './Button';

const imgTransparent =
  'url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAALEgAACxIB0t1+/AAAABx0RVh0U29mdHdhcmUAQWRvYmUgRmlyZXdvcmtzIENTNXG14zYAAAAWdEVYdENyZWF0aW9uIFRpbWUAMDIvMDMvMTSFgaz/AAAADUlEQVQImWP4//+/LwAJSQNLpIypVwAAAABJRU5ErkJggg==)';

export function Modal() {
  return (
    <div
      className='flex justify-center items-center'
      style={{
        position: 'absolute',
        top: 0,
        width: '100%',
        height: '100%',
        background: imgTransparent,
      }}
    >
      <div
        className=''
        style={{
          opacity: 1,
        }}
      >
        <h4 className='text-center'>
          Olá username, Queremos lhe desejar boa sorte e boas-vindas, que você não desista do jogo
          antes do G Rank!
        </h4>
        <div className='flex items-center justify-center flex-col'>
          <Button>Entrar</Button>
          <Button>Trocar de conta</Button>
        </div>
      </div>
    </div>
  );
}
