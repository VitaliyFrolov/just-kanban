import { getRequestConfig } from 'next-intl/server';

export const SUPPORTED_LOCALES = ['en', 'ru'] as const;
export const LOCALE_PREFIX = 'always' as const;

export default getRequestConfig(async ({ locale, requestLocale }) => {
  const reqLocale = await requestLocale;
  const validLocale =
    !!reqLocale &&
    SUPPORTED_LOCALES.includes(reqLocale as (typeof SUPPORTED_LOCALES)[number]);
  const activeLocale = validLocale
    ? (reqLocale as (typeof SUPPORTED_LOCALES)[number])
    : locale;

  return {
    locale: activeLocale ?? SUPPORTED_LOCALES[0],
  };
});
