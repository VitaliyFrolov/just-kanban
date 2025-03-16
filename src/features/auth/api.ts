'use server';

import { ApiPath } from 'shared/config';
import { customFetch } from 'shared/api';

interface RegistrationData {
  username: string;
  email: string;
  password: string;
}

export const registration = async ({
  username,
  email,
  password,
}: RegistrationData) =>
  customFetch.POST<string>(ApiPath.Registration, {
    body: JSON.stringify({ username, email, password }),
  });

export const login = async (identifier: string, password: string) =>
  customFetch.POST<string>(ApiPath.Login, {
    body: JSON.stringify({ identifier, password }),
  });
