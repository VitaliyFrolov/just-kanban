'use server';

import { ApiPath } from 'shared/config';
import { ApiValidationError } from 'shared/lib';
import { customFetch } from 'shared/api';

interface RegistrationData {
  username: string;
  email: string;
  password: string;
  firstName: string;
  lastName: string;
}

export const registration = async ({
  username,
  email,
  password,
  firstName,
  lastName,
}: RegistrationData) =>
  customFetch.POST<
    string,
    ApiValidationError<{
      username: string;
      email: string;
      password: string;
      first_name: string;
      last_name: string;
    }>
  >(ApiPath.Registration, {
    body: JSON.stringify({
      username,
      email,
      password,
      first_name: firstName,
      last_name: lastName,
    }),
  });

export const login = async (identifier: string, password: string) =>
  customFetch.POST<string, ApiValidationError<{ root: string }>>(
    ApiPath.Login,
    {
      body: JSON.stringify({ identifier, password }),
    },
  );
