'use server';

import { cookies } from 'next/headers';
import { jwtDecode } from 'jwt-decode';

import { ACCESS_TOKEN_COOKIE_NAME } from 'shared/config';

export interface AuthUserClaims {
  id: string;
  email: string;
  avatar: string;
  username: string;
  first_name: string;
  last_name: string;
}

export const getUserClaims = async (token: string): Promise<AuthUserClaims> => {
  return jwtDecode<AuthUserClaims>(token);
};

export const getAccessToken = async () => {
  const cookie = await cookies();
  cookie.get(ACCESS_TOKEN_COOKIE_NAME);
};

export const setAccessToken = async (token: string) => {
  const cookie = await cookies();
  cookie.set(ACCESS_TOKEN_COOKIE_NAME, token);
};

export const removeAccessToken = async () => {
  const cookie = await cookies();
  cookie.delete(ACCESS_TOKEN_COOKIE_NAME);
};
