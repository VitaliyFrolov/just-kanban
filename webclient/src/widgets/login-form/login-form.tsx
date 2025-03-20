'use client';

import { AppConfig } from 'next-intl';
import { FC } from 'react';

import { Button, InformationMessage, TextField } from 'shared/ui';
import { clsx, formResolver, useForm, useRouter, validate } from 'shared/lib';
import { RoutePath } from 'shared/config';
import { login } from 'features/auth';
import { setAccessToken } from 'features/auth';

interface LoginFormProps {
  className?: string;
  messages?: {
    auth?: AppConfig['Messages']['auth'];
    validation?: AppConfig['Messages']['validation'];
  };
}

export const LoginForm: FC<LoginFormProps> = ({ className, messages }) => {
  const router = useRouter();
  const schema = validate.object({
    identifier: validate
      .string()
      .nonempty(messages?.validation?.['validation-required']),
    password: validate
      .string()
      .nonempty(messages?.validation?.['validation-required']),
    remember: validate.boolean(),
  });
  const {
    register,
    handleSubmit,
    setError,
    formState: { isSubmitting, errors },
  } = useForm({
    defaultValues: {
      remember: false,
    },
    resolver: formResolver(schema),
  });

  const onSubmit = handleSubmit(async (fields) => {
    const response = await login(fields.identifier, fields.password);
    if (response.ok) {
      await setAccessToken(response.data);
      router.push(RoutePath.Home);
    } else {
      if (response.error.fields?.root) {
        setError('root', { message: response.error.fields.root });
      }
    }
  });

  return (
    <div
      className={clsx(
        className,
        'flex flex-col gap-xl bg-surface p-l rounded-m shadow-surface tablet:max-w-[630px]',
      )}
    >
      <h4 className="text-xl font-bold whitespace-pre-wrap tablet:text-xxl">
        {messages?.auth?.['sign-in-tagline']}
      </h4>
      <form
        className="flex flex-col gap-l"
        onSubmit={onSubmit}
      >
        <TextField
          {...register('identifier')}
          autoComplete="username webauthn"
          placeholder={messages?.auth?.['identifier-input']}
          errorMessage={errors.identifier?.message}
        />
        <TextField
          {...register('password')}
          type="password"
          autoComplete="current-password"
          placeholder={messages?.auth?.['password-input']}
          errorMessage={errors.password?.message}
        />
        {errors.root && (
          <InformationMessage error>{errors.root.message}</InformationMessage>
        )}
        <Button
          type="submit"
          disabled={isSubmitting}
          loading={isSubmitting}
          variant="primaryFilled"
        >
          {messages?.auth?.['login-button']}
        </Button>
      </form>
    </div>
  );
};
