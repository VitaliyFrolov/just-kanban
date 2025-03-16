import { FC, PropsWithChildren } from 'react';
import { Metadata } from 'next';
import { PT_Sans } from 'next/font/google';
import { PT_Serif } from 'next/font/google';
import { cookies } from 'next/headers';

import 'app/styles/index.css';
import { IntlProvider, ThemeProvider } from 'app/providers';
import { THEME_ATTRIBUTE_NAME, THEME_COOKIE_NAME, Theme } from 'shared/config';
import { RoutingDefaultParams } from 'shared/types';
import { clsx } from 'shared/lib';

type RootLayoutProps = PropsWithChildren<{
  params: Promise<RoutingDefaultParams>;
}>;

export const generateMetadata = async (): Promise<Metadata> => {
  return {
    title: '',
    description: '',
  };
};

const PTSans = PT_Sans({
  display: 'swap',
  preload: true,
  variable: '--font-sans',
  weight: ['700', '400'],
});

const PTSerif = PT_Serif({
  display: 'swap',
  variable: '--font-serif',
  preload: true,
  weight: ['400', '700'],
});

const Providers: FC<PropsWithChildren> = async (props) => {
  return (
    <IntlProvider>
      <ThemeProvider>{props.children}</ThemeProvider>
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
        <Providers>{children}</Providers>
      </body>
    </html>
  );
};

export default RootLayout;
