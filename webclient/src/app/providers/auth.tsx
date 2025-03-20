'use client';

import { FC, PropsWithChildren } from 'react';

import { AuthContext, AuthUserClaims } from 'features/auth';

interface AuthProviderProps extends PropsWithChildren {
  user: AuthUserClaims | null;
}

export const AuthProvider: FC<AuthProviderProps> = ({ user, children }) => {
  return <AuthContext value={{ user }}>{children}</AuthContext>;
};
