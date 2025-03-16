import { ACCESS_TOKEN_HEADER_NAME, API_URL } from 'shared/config';
import { cleanObject } from 'shared/lib';

import { readAuthToken } from './read-auth-token';

export interface ListFetchResponse<T> {
  count: number;
  next: number | null;
  previous: number | null;
  results: T[];
}

interface CustomFetchOptions extends RequestInit {
  query?: Record<string, unknown>;
  // Defines if origin must be taken from provided url
  customOrigin?: boolean;
}

export interface HttpError<TBody = unknown> extends Error {
  status: number;
  body: TBody;
}

export const customFetch = {
  async _createRequest<TResponse = unknown, TErrorBody = unknown>(
    url: string,
    options?: CustomFetchOptions,
  ) {
    const authToken = await readAuthToken();

    const headers = new Headers({
      ...options?.headers,
      ...(!(options?.body instanceof FormData) && {
        'Content-Type': 'application/json',
      }),
    });

    if (authToken) {
      headers.set(ACCESS_TOKEN_HEADER_NAME, `Bearer ${authToken}`);
    }

    const cleanedQuery = options?.query ? cleanObject(options.query) : {};

    const searchParams = new URLSearchParams(
      cleanedQuery as Record<string, string>,
    ).toString();

    console.log(API_URL, url);

    const res = await fetch(
      `${options?.customOrigin ? '' : API_URL}${url}${
        searchParams && '?' + searchParams
      }`,
      {
        ...options,
        headers,
      },
    );

    const bodyText = (await res.text()) || JSON.stringify(null);

    if (!res.ok) {
      throw Object.assign(
        new Error(`Failed to fetch ${url}: ${res.status} ${res.statusText}`),
        {
          status: res.status,
          body: JSON.parse(bodyText),
        },
      ) as HttpError<TErrorBody>;
    }

    return JSON.parse(bodyText) as TResponse;
  },

  POST<TResponse>(url: string, options?: CustomFetchOptions) {
    return customFetch._createRequest<TResponse>(url, {
      ...options,
      method: 'POST',
    });
  },

  GET<TResponse>(url: string, options?: CustomFetchOptions) {
    return customFetch._createRequest<TResponse>(url, {
      ...options,
      method: 'GET',
    });
  },

  PATCH<TResponse>(url: string, options?: CustomFetchOptions) {
    return customFetch._createRequest<TResponse>(url, {
      ...options,
      method: 'PATCH',
    });
  },

  DELETE<TResponse>(url: string, options?: CustomFetchOptions) {
    return customFetch._createRequest<TResponse>(url, {
      ...options,
      method: 'DELETE',
    });
  },
};
