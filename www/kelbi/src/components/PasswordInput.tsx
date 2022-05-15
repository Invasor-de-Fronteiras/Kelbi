// import { useField } from 'formik';
import React from 'react';

// interface {}
export function PasswordInput() {
  // const a = useField()
  return (
    <div className='input-group'>
      {/* <label>Password</label> */}
      <input className='input' placeholder='Password' required />
    </div>
  );
}
