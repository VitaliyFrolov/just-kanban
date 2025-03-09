'use server';

import { cookies } from 'next/headers';

import { ACCESS_TOKEN_COOKIE_NAME } from 'shared/config/auth';

export const readAuthToken = async () => {
  const cookie = await cookies();
  return cookie.get(ACCESS_TOKEN_COOKIE_NAME)?.value;
};
