'use client';

import { createContext } from 'react';

import { Theme } from 'shared/config';

interface ThemeContextState {
  theme: Theme;
  setTheme: (newTheme: Theme) => void;
}

export const ThemeContext = createContext<ThemeContextState>({
  theme: Theme.default,
  setTheme: () => null,
});
