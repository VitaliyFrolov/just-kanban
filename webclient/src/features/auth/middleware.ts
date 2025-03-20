'use server';

import { NextRequest } from 'next/server';
import { cookies } from 'next/headers';

import {
  ACCESS_TOKEN_COOKIE_NAME,
  RoutePath,
  SECURE_ROUTES_PATHS,
} from 'shared/config';
import { getUserClaims } from 'features/auth';

export const authMiddleware = async (req: NextRequest) => {
  if (
    !SECURE_ROUTES_PATHS.some(
      (routePath) =>
        routePath.includes(req.nextUrl.pathname) ||
        req.nextUrl.pathname.includes(routePath),
    )
  ) {
    return req;
  }

  const cookie = await cookies();
  const accessToken = cookie.get(ACCESS_TOKEN_COOKIE_NAME)?.value;
  const userClaims = accessToken ? await getUserClaims(accessToken) : null;

  if (userClaims) {
    return req;
  } else {
    req.nextUrl.pathname = RoutePath.Login;
  }

  return req;
};
