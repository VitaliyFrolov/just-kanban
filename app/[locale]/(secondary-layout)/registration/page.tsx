import { FC } from 'react';
import { Metadata } from 'next';

import { RegistrationPage } from 'pages/registration';
import { getMessages } from 'shared/lib';

export const generateMetadata = async (): Promise<Metadata> => {
  const messages = await getMessages();

  return {
    title: messages.auth.registration,
  };
};

const RoutingPage: FC = () => {
  return <RegistrationPage />;
};

export default RoutingPage;
