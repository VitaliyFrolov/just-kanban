'use client';

import { AppConfig, NextIntlClientProvider } from 'next-intl';
import { FC, PropsWithChildren } from 'react';

export const IntlProvider: FC<
  PropsWithChildren & { locale: AppConfig['Locale'] }
> = ({ children, locale }) => {
  return (
    <NextIntlClientProvider
      messages={null}
      locale={locale}
    >
      {children}
    </NextIntlClientProvider>
  );
};
