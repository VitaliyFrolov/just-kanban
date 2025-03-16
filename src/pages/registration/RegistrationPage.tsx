import { FC } from 'react';

import { RegistrationForm } from 'widgets/registration-form';
import { getMessages } from 'shared/lib';

export const RegistrationPage: FC = async () => {
  const messages = await getMessages();

  return (
    <RegistrationForm
      messages={{ auth: messages.auth, validation: messages.validation }}
    />
  );
};
