import React from 'react';
import BarbaryWalking from '../assets/loading/banbaro-walking.gif';

interface LauncherLoadingProps {
  message: string;
  showUpdateProgress: boolean;
  fileProgress: number;
  totalProgress: number;
}

export function LauncherLoading({
  message,
  showUpdateProgress,
  fileProgress,
  totalProgress,
}: LauncherLoadingProps) {
  return (
    <div className='flex items-center justify-center flex-col'>
      <img id='img-bg' src={BarbaryWalking} height='100%' width='100%' />
      <span
        style={{
          marginTop: '5px',
        }}
      >
        {message}
      </span>
      {showUpdateProgress && (
        <div className='w-full mt-6'>
          <div className='w-full bg-gray-200 rounded-full h-2 mb-4'>
            <div
              className='bg-orange-500 h-2 rounded-full'
              style={{ width: `${fileProgress}%` }}
            ></div>
          </div>
          <div className='w-full bg-gray-200 rounded-full h-2 mb-4'>
            <div
              className='bg-green-600 h-2 rounded-full'
              style={{ width: `${totalProgress}%` }}
            ></div>
          </div>
        </div>
      )}
    </div>
  );
}
