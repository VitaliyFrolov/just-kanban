import { FC } from 'react';

import { LoginForm } from 'widgets/login-form';
import { getMessages } from 'shared/lib';

export const LoginPage: FC = async () => {
  const messages = await getMessages();

  return (
    <LoginForm
      messages={{ auth: messages.auth, validation: messages.validation }}
    />
  );
};
