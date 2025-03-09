import { FC, PropsWithChildren } from 'react';
import { Metadata } from 'next';
import { cookies } from 'next/headers';
import { getMessages } from 'next-intl/server';

import 'app/styles/index.css';
import { IntlProvider, ThemeProvider } from 'app/providers';
import { THEME_ATTRIBUTE_NAME, THEME_COOKIE_NAME, Theme } from 'shared/config';

type RootLayoutProps = PropsWithChildren<{
  params: {
    locale: string;
  };
}>;

export const generateMetadata = async (): Promise<Metadata> => {
  return {
    title: '',
    description: '',
  };
};

const Providers: FC<PropsWithChildren> = async (props) => {
  const [messages] = await Promise.all([getMessages()]);

  return (
    <IntlProvider messages={messages}>
      <ThemeProvider>{props.children}</ThemeProvider>
    </IntlProvider>
  );
};

const RootLayout: FC<RootLayoutProps> = async ({ params, children }) => {
  const cookie = await cookies();
  const theme = cookie.get(THEME_COOKIE_NAME)?.value ?? Theme.default;

  const htmlAttributes = {
    lang: params.locale,
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
