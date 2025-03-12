import { FC, PropsWithChildren } from 'react';
import { Metadata } from 'next';
import { cookies } from 'next/headers';

import 'app/styles/index.css';
import { IntlProvider, ThemeProvider } from 'app/providers';
import { THEME_ATTRIBUTE_NAME, THEME_COOKIE_NAME, Theme } from 'shared/config';
import { RoutingDefaultParams } from 'shared/types';

type RootLayoutProps = PropsWithChildren<{
  params: Promise<RoutingDefaultParams>;
}>;

export const generateMetadata = async (): Promise<Metadata> => {
  return {
    title: '',
    description: '',
  };
};

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
      <body className="bg-primary text-secondary">
        <Providers>{children}</Providers>
      </body>
    </html>
  );
};

export default RootLayout;
