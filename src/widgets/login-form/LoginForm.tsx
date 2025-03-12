'use client';

import { FC } from 'react';

import { Button, TextField } from 'shared/ui';
import { formResolver, useForm, useMessages, validate } from 'shared/lib';

interface LoginFormProps {
  className?: string;
}

export const LoginForm: FC<LoginFormProps> = ({ className }) => {
  const { Provider, loading, t } = useMessages(['auth']);
  const schema = validate.object({
    username: validate.string(),
    password: validate.string(),
    remember: validate.boolean(),
  });
  const {
    register,
    handleSubmit,
    formState: { isSubmitting, isDirty, isValid },
  } = useForm({
    resolver: formResolver(schema),
  });

  const onSubmit = handleSubmit(async () => {});

  return (
    <Provider>
      {loading ? (
        'Loading'
      ) : (
        <form
          className={className}
          onSubmit={onSubmit}
        >
          <TextField
            {...register('username')}
            autoComplete="username"
            placeholder={t('auth.identifier-input')}
          />
          <TextField
            {...register('password')}
            type="password"
            autoComplete="current-password"
            placeholder={t('auth.password-input')}
          />
          <Button
            type="submit"
            disabled={!isValid || !isDirty || isSubmitting}
          >
            {t('auth.login-button')}
          </Button>
        </form>
      )}
    </Provider>
  );
};
