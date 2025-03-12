'use client';

import { AppConfig, NextIntlClientProvider, useLocale } from 'next-intl';
import { PropsWithChildren, ReactNode, useEffect, useState } from 'react';

export function useMessages<NSName extends keyof AppConfig['Messages']>(
  namespaces: NSName[],
) {
  const locale = useLocale();
  const [error, setError] = useState<Error | null>(null);
  const [loading, setLoading] = useState(true);
  const [messages, setMessages] = useState<Record<
    NSName,
    AppConfig['Messages'][NSName]
  > | null>(null);

  useEffect(() => {
    setLoading(true);
    Promise.all(
      namespaces.map((namespace) =>
        fetch(`/messages/${locale}/${namespace}.json`)
          .then((res) => {
            if (!res.ok)
              throw new Error(`Failed to load i18n namespace ${namespace}`);
            return res.json();
          })
          .then((data) => ({ [namespace]: data })),
      ),
    )
      .then((results) => {
        setMessages(Object.assign({}, ...results));
        setLoading(false);
      })
      .catch((err) => {
        setError(err);
        setLoading(false);
      });
  }, [locale, ...namespaces]);

  type TranslationKey = {
    [N in NSName]: `${N}.${keyof AppConfig['Messages'][N] & string}`;
  }[NSName];

  const t = (
    key: TranslationKey,
    variables?: Record<string, string | number>,
  ): string => {
    if (loading || error || !messages) return key;
    const [namespace, translationKey] = key.split('.') as [
      NSName,
      keyof AppConfig['Messages'][NSName],
    ];
    const msg = messages[namespace];
    if (!msg) return key;
    const translation = msg[translationKey];
    if (!translation) return key;
    if (variables) {
      return Object.entries(variables).reduce(
        (str, [varName, value]) => str.replace(`{${varName}}`, String(value)),
        translation as string,
      );
    }
    return translation as string;
  };

  const tWithRich = Object.assign(t, {
    rich: (
      key: TranslationKey,
      options?: Record<
        string,
        string | number | ((chunk: string) => ReactNode)
      >,
    ): string | ReactNode => {
      const translation = t(key);
      if (!options) return translation;
      let result = translation;
      for (const [tag, value] of Object.entries(options)) {
        if (typeof value === 'function') {
          result = result.replace(`{${tag}}`, String(value(tag)));
        } else {
          result = result.replace(`{${tag}}`, String(value));
        }
      }
      return result;
    },
  });

  const LazyTranslationProvider = ({ children }: PropsWithChildren) => {
    if (loading || error || !messages) return <>{children}</>;
    return (
      <NextIntlClientProvider
        locale={locale}
        messages={messages}
      >
        {children}
      </NextIntlClientProvider>
    );
  };

  return { t: tWithRich, Provider: LazyTranslationProvider, loading, error };
}
