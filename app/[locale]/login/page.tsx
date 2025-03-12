import { FC } from 'react';
import { Metadata } from 'next';

import { LoginPage } from 'pages/login';
import { getMessages } from 'shared/lib';

export const generateMetadata = async (): Promise<Metadata> => {
  const t = await getMessages('auth');

  return {
    title: t('page-title'),
  };
};

const RoutingPage: FC = () => {
  return <LoginPage />;
};

export default RoutingPage;
