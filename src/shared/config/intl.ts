import { getRequestConfig } from 'next-intl/server';

export const SUPPORTED_LOCALES = ['en', 'ru'] as const;
export const LOCALE_PREFIX = 'always' as const;

export default getRequestConfig(async ({ locale, requestLocale }) => {
  const reqLocale = await requestLocale;
  const isValidLocale =
    !!reqLocale &&
    SUPPORTED_LOCALES.includes(reqLocale as (typeof SUPPORTED_LOCALES)[number]);
  const activeLocale = isValidLocale
    ? (reqLocale as (typeof SUPPORTED_LOCALES)[number])
    : locale;

  return {
    locale: activeLocale ?? SUPPORTED_LOCALES[0],
    messages: {
      common: {
        ...(await import(`public/messages/${activeLocale}/common.json`)),
      },
      auth: {
        ...(await import(`public/messages/${activeLocale}/auth.json`)),
      },
      validation: {
        ...(await import(`public/messages/${activeLocale}/validation.json`)),
      },
    },
  };
});
