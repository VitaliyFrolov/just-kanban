import { FC, PropsWithChildren } from 'react';
import { PT_Sans, PT_Serif } from 'next/font/google';
import { AppConfig } from 'next-intl';
import { cookies } from 'next/headers';

import 'app/styles/index.css';
import {
  ACCESS_TOKEN_COOKIE_NAME,
  THEME_ATTRIBUTE_NAME,
  THEME_COOKIE_NAME,
  Theme,
} from 'shared/config';
import { AuthProvider, IntlProvider, ThemeProvider } from 'app/providers';
import { RoutingDefaultParams } from 'shared/types';
import { clsx } from 'shared/lib';
import { getUserClaims } from 'features/auth';

type RootLayoutProps = PropsWithChildren<{
  params: Promise<RoutingDefaultParams>;
}>;

const PTSans = PT_Sans({
  display: 'swap',
  preload: true,
  variable: '--font-sans',
  weight: ['400', '700'],
});

const PTSerif = PT_Serif({
  display: 'swap',
  preload: true,
  variable: '--font-serif',
  weight: ['400', '700'],
});

const Providers: FC<
  PropsWithChildren & { locale: AppConfig['Locale'] }
> = async ({ locale, children }) => {
  const cookie = await cookies();
  const jwtToken = cookie.get(ACCESS_TOKEN_COOKIE_NAME)?.value;
  const userClaims = jwtToken ? getUserClaims(jwtToken) : null;

  return (
    <IntlProvider locale={locale}>
      <AuthProvider user={userClaims}>
        <ThemeProvider>{children}</ThemeProvider>
      </AuthProvider>
    </IntlProvider>
  );
};

const RootLayout: FC<RootLayoutProps> = async ({ params, children }) => {
  const { locale } = await params;
  const cookie = await cookies();
  const theme = cookie.get(THEME_COOKIE_NAME)?.value ?? Theme.default;

  const htmlAttributes = {
    lang: locale,
    [THEME_ATTRIBUTE_NAME]: theme,
  };

  return (
    <html {...htmlAttributes}>
      <body
        className={clsx(
          'bg-site-background',
          PTSans.variable,
          PTSerif.variable,
        )}
      >
        <Providers locale={locale}>{children}</Providers>
      </body>
    </html>
  );
};

export default RootLayout;
