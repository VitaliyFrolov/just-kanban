'use client';

import { createContext } from 'react';

import { AuthUserClaims } from './jwt';

interface AuthContextState {
  user: AuthUserClaims | null;
}

export const AuthContext = createContext<AuthContextState>({
  user: null,
});
