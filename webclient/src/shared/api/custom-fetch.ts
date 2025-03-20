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
}

type CustomFetchResponse<TData = unknown, TErrorBody = unknown> =
  | {
      ok: true;
      data: TData;
      error?: never;
    }
  | {
      ok: false;
      data?: never;
      error: Partial<TErrorBody>;
    };

export const customFetch = {
  async _createRequest<TResponse = unknown, TErrorBody = unknown>(
    url: string,
    options?: CustomFetchOptions,
  ): Promise<CustomFetchResponse<TResponse, TErrorBody>> {
    const authToken = await readAuthToken();
    const isCustomOrigin = url.startsWith('https://');

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
      `${isCustomOrigin ? '' : API_URL}${url}${(url.includes('?') ? '&' : '?') + searchParams}`,
      {
        ...options,
        headers,
      },
    );

    const bodyText = (await res.text()) || JSON.stringify(null);

    if (!res.ok) {
      return {
        ok: false,
        error: JSON.parse(bodyText),
      };
    }

    return {
      ok: true,
      data: JSON.parse(bodyText),
    };
  },

  POST<TResponse, TError = unknown>(url: string, options?: CustomFetchOptions) {
    return customFetch._createRequest<TResponse, TError>(url, {
      ...options,
      method: 'POST',
    });
  },

  GET<TResponse, TError = unknown>(url: string, options?: CustomFetchOptions) {
    return customFetch._createRequest<TResponse, TError>(url, {
      ...options,
      method: 'GET',
    });
  },

  PATCH<TResponse, TError = unknown>(
    url: string,
    options?: CustomFetchOptions,
  ) {
    return customFetch._createRequest<TResponse, TError>(url, {
      ...options,
      method: 'PATCH',
    });
  },

  DELETE<TResponse, TError = unknown>(
    url: string,
    options?: CustomFetchOptions,
  ) {
    return customFetch._createRequest<TResponse, TError>(url, {
      ...options,
      method: 'DELETE',
    });
  },
};
