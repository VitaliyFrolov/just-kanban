'use client';

import { createContext } from 'react';

interface FormControlContextState {
  name?: string;
  required?: boolean;
  error?: boolean;
  disabled?: boolean;
}

export const FormControlContext = createContext<FormControlContextState>({
  required: false,
  error: false,
  disabled: false,
});
