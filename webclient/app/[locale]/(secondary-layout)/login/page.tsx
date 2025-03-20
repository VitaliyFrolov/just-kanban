import { FC } from 'react';
import { Metadata } from 'next';

import { LoginPage } from 'pages/login';
import { getMessages } from 'shared/lib';

export const generateMetadata = async (): Promise<Metadata> => {
  const messages = await getMessages();

  return {
    title: messages.auth['sign-in'],
  };
};

const RoutingPage: FC = () => {
  return <LoginPage />;
};

export default RoutingPage;
