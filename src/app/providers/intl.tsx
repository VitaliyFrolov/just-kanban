import { FC, PropsWithChildren } from 'react';
import { NextIntlClientProvider } from 'next-intl';

export const IntlProvider: FC<PropsWithChildren> = ({ children }) => {
  return (
    <NextIntlClientProvider messages={null}>{children}</NextIntlClientProvider>
  );
};
